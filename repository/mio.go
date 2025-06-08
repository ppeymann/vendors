package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/utils"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type mioRepo struct {
	pg       *gorm.DB
	minio    *minio.Client
	opts     config.StorageOptions
	database string
	table    string
	secret   string
}

// GetObject implements models.MioRepository.
func (r *mioRepo) GetObject(in *models.DownloadInput) ([]byte, *models.StorageEntity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	aid, err := strconv.Atoi(in.AccountId)
	if err != nil {
		return nil, nil, utils.ErrUserPrincipalsNotFound
	}

	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, nil, utils.ErrUserPrincipalsNotFound
	}

	file := &models.StorageEntity{}

	query := "account = ? AND tag = ? AND id = ? AND bucket = ?"
	bn := fmt.Sprintf("user-%s-bucket", in.AccountId)

	res := r.pg.Where(query, aid, in.Tag, id, bn).First(file)
	if res.Error != nil {
		return nil, nil, err
	}

	object, err := r.minio.GetObject(ctx, bn, file.FileName, minio.GetObjectOptions{
		ServerSideEncryption: nil,
	})

	if err != nil {
		return nil, nil, err
	}

	defer func(object *minio.Object) {
		_ = object.Close()
	}(object)

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, nil, err
	}

	return data, file, nil
}

// GetResizedImageObject implements models.MioRepository.
func (r *mioRepo) GetResizedImageObject(in *models.DownloadInput) ([]byte, *models.StorageEntity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	aid, err := strconv.Atoi(in.AccountId)
	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, nil, utils.ErrUserPrincipalsNotFound
	}

	file := &models.StorageEntity{}

	query := "account = ? AND tag = ? AND id = ? AND bucket = ?"
	bn := fmt.Sprintf("user-%s-bucket", in.AccountId)

	res := r.pg.Where(query, aid, in.Tag, id, bn).First(file)
	if res.Error != nil {
		return nil, nil, err
	}

	fn := fmt.Sprintf("_%d.jpg", in.Size)

	sizedFn := strings.Replace(file.FileName, ".jpg", fn, -1)

	_, err = r.minio.StatObject(ctx, bn, sizedFn, minio.StatObjectOptions{})

	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {

			object, e := r.minio.GetObject(ctx, bn, file.FileName, minio.GetObjectOptions{})
			if e != nil {
				return nil, nil, vendora.ErrNotFound
			}

			defer func(object *minio.Object) {
				_ = object.Close()
			}(object)

			data, e := io.ReadAll(object)
			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			reader := bytes.NewReader(data)
			img, e := imaging.Decode(reader, imaging.AutoOrientation(false))
			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			resized := imaging.Resize(img, int(in.Size), 0, imaging.Lanczos)

			buff := new(bytes.Buffer)
			e = imaging.Encode(buff, resized, imaging.JPEG, imaging.JPEGQuality(100))
			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			temp, e := os.MkdirTemp("./data", bn)
			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			path := fmt.Sprintf("%s/%s", temp, sizedFn)

			e = imaging.Save(resized, path, imaging.JPEGQuality(100))
			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			defer func() {
				_ = os.RemoveAll(temp)
			}()

			_, e = r.minio.FPutObject(ctx, bn, sizedFn, path, minio.PutObjectOptions{
				UserTags:    map[string]string{"target": file.Tag},
				ContentType: file.ContentType,
			})

			if e != nil {
				return nil, nil, vendora.ErrInternalServer
			}

			return buff.Bytes(), file, nil
		}
	} else {
		object, err := r.minio.GetObject(ctx, bn, sizedFn, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, err
		}

		defer func(object *minio.Object) {
			_ = object.Close()
		}(object)

		data, err := io.ReadAll(object)
		if err != nil {
			return nil, nil, err
		}

		return data, file, nil
	}

	return nil, nil, nil

}

// PutObject implements models.MioRepository.
func (r *mioRepo) PutObject(bucketName string, objectName string, path string, ct string, ext string, tag models.ObjectTag) (*vendora.BaseResult, error) {
	aid, err := strconv.Atoi(bucketName)
	if err != nil {
		return nil, utils.ErrUserPrincipalsNotFound
	}

	bn := fmt.Sprintf("user-%s-bucket", bucketName)
	fileId := ksuid.New().String()

	file := &models.StorageEntity{
		Model:       gorm.Model{},
		Account:     uint(aid),
		ContentType: ct,
		Tag:         objectName,
		Bucket:      bn,
		FileName:    fmt.Sprintf("%s%s", fileId, ext),
	}

	err = r.pg.Transaction(func(tx *gorm.DB) error {
		if res := r.Model().Create(file).Error; res != nil {
			str := res.(*pgconn.PgError).Message
			if strings.Contains(str, "duplicate key value") {
				return models.ErrAccountExist
			}
			return res
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	data := fmt.Sprintf(`{"account_id":"%d","tag":"%s","id":"%d"}`, aid, objectName, file.ID)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	err = r.minio.MakeBucket(ctx, bn, minio.MakeBucketOptions{Region: r.opts.Region})
	if err != nil {
		exists, existErr := r.minio.BucketExists(ctx, bn)
		if existErr != nil && !exists {
			return nil, err
		}
	}

	info, err := r.minio.FPutObject(ctx, bn, file.FileName, path, minio.PutObjectOptions{
		UserTags:    map[string]string{"target": string(tag)},
		ContentType: ct,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(info)

	cypher, err := utils.EncryptText(data, r.secret)
	if err != nil {
		return nil, vendora.ErrInternalServer
	}

	return &vendora.BaseResult{
		Result: cypher,
	}, err
}

// RemoveObject implements models.MioRepository.
func (r *mioRepo) RemoveObject(bucketName string, objectName string, path string, objectId string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := r.minio.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	err = r.minio.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Migrate implements models.MioRepository.
func (r *mioRepo) Migrate() error {
	return r.pg.AutoMigrate(&models.StorageEntity{})
}

// Model implements models.MioRepository.
func (r *mioRepo) Model() *gorm.DB {
	return r.pg.Model(&models.StorageEntity{})
}

// Name implements models.MioRepository.
func (r *mioRepo) Name() string {
	return r.table
}

func NewMioRepo(pg *gorm.DB, opts config.StorageOptions, database string) (models.MioRepository, error) {
	client, err := minio.New(opts.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.User, opts.Secret, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}

	return &mioRepo{
		pg:       pg,
		minio:    client,
		opts:     opts,
		database: database,
		table:    "mio_entities",
		secret:   env.GetEnv("JWT", ""),
	}, nil
}
