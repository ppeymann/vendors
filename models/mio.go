package models

import (
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/auth"
	"gorm.io/gorm"
)

type (
	ObjectTag string

	// MioService represents method signatures for api mio service endpoint.
	// so any object that stratifying this interface can be used as mio service for api endpoint.
	MioService interface {
		Upload(in *UploadInput, ctx *gin.Context) *vendora.BaseResult
		Download(in *DownloadInput, ctx *gin.Context) ([]byte, *StorageEntity, error)
		Image(in *DownloadInput, ctx *gin.Context) ([]byte, *StorageEntity, error)
	}

	// MioRepository represents method signatures for mio service domain repository.
	// so any object that stratifying this interface can be used as mio service domain repository
	MioRepository interface {
		PutObject(bucketName, objectName, path, ct, ext string, tag ObjectTag) (*vendora.BaseResult, error)
		GetObject(in *DownloadInput) ([]byte, *StorageEntity, error)
		GetResizedImageObject(in *DownloadInput) ([]byte, *StorageEntity, error)
		RemoveObject(bucketName, objectName, path, objectId string) error

		vendora.BaseRepository
	}

	// MioHandler represents method signatures for api mio service handler endpoint.
	// so any object that stratifying this interface can be used as mio service handler endpoint.
	MioHandler interface {
		Upload(ctx *gin.Context)
		Download(ctx *gin.Context)
		Image(ctx *gin.Context)
	}

	// UploadInput is DTO for transferring file upload request params.
	UploadInput struct {
		Claims      *auth.Claims `json:"-"`
		Tag         string       `json:"tag"`
		Size        int64        `json:"size"`
		ContentType string       `json:"content_type"`
		FileName    string       `json:"file_name"`
	}

	// DownloadInput is DTO for transferring file download request params.
	DownloadInput struct {
		Token     string `json:"-"`
		AccountId string `json:"account_id"`
		Tag       string `json:"tag"`
		Id        string `json:"id"`
		Size      uint   `json:"-"`
	}

	// StorageEntity is entity of storage file object item info.
	//
	// swagger:model StorageEntity
	StorageEntity struct {
		gorm.Model

		// Account id of object owner
		Account uint `json:"account" gorm:"column:account;index"`

		// Tag for stored object to manage access permission
		Tag string `json:"tag" gorm:"column:tag;index"`

		// Bucket name of stored object that
		Bucket string `json:"bucket" gorm:"column:bucket;index"`

		// ContentType of stored object
		ContentType string `json:"content_type" gorm:"column:content_type;index"`

		// ContentType of stored object for put back on download request
		FileName string `json:"file_name" gorm:"column:file_name;index"`
	}
)

const (
	PublicTag  ObjectTag = "public"
	PrivateTag ObjectTag = "private"
	ChatTag    ObjectTag = "chat"
	ProfileTag ObjectTag = "profile"
)
