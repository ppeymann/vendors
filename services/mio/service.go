package mio

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
)

type service struct {
	opts config.StorageOptions
	repo models.MioRepository
}

// Download implements models.MioService.
func (s *service) Download(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	return s.repo.GetObject(in)
}

// Image implements models.MioService.
func (s *service) Image(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	data, file, err := s.repo.GetResizedImageObject(in)
	if err != nil {
		return nil, nil, err
	}

	return data, file, nil
}

// Upload implements models.MioService.
func (s *service) Upload(in *models.UploadInput, ctx *gin.Context) *vendora.BaseResult {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{ErrBadUploadRequest.Error()},
			Status: http.StatusOK,
		}
	}

	tempName := fmt.Sprintf("upload_%d_%s_%s", in.Claims.Subject, in.Tag, formFile.Filename)

	fp := filepath.Clean("./data/" + tempName)
	resFile, err := os.Create(fp)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	fex := filepath.Ext(resFile.Name())

	file, err := formFile.Open()
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	_, err = io.Copy(resFile, file)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	defer func() {
		_ = resFile.Close()
		_ = os.Remove(resFile.Name())
	}()

	fn := resFile.Name()
	result, err := s.repo.PutObject(strconv.Itoa(int(in.Claims.Subject)), in.Tag, fn, in.ContentType, fex, models.ObjectTag(in.Tag))
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return result
}

func NewService(opts config.StorageOptions, repo models.MioRepository) models.MioService {
	return &service{
		opts: opts,
		repo: repo,
	}
}
