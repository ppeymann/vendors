package mio

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/utils"
	validations "github.com/ppeymann/vendors.git/validation"
)

// ErrValidationFailed is returned request json body object is not stratified schema rules
var ErrValidationFailed = errors.New("request rejected, validation error")
var ErrBadUploadRequest = errors.New("upload request is illegal")
var ErrFileTypeNotSupported = errors.New("upload file type is not supported")
var ErrInvalidAccessToken = errors.New("invalid photo access token")

type validationService struct {
	schemas map[string][]byte
	secret  string
	opts    config.StorageOptions
	next    models.MioService
}

// Download implements models.MioService.
func (v *validationService) Download(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	plain, err := utils.DecryptText(in.Token, v.secret)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal([]byte(plain), in)
	if err != nil {
		return nil, nil, ErrInvalidAccessToken
	}

	validationErr := validations.Validate(in, v.schemas)
	if validationErr != nil {
		return nil, nil, ErrValidationFailed
	}

	return v.next.Download(in, ctx)
}

// Image implements models.MioService.
func (v *validationService) Image(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	plain, err := utils.DecryptText(in.Token, v.secret)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal([]byte(plain), in)
	if err != nil {
		return nil, nil, ErrInvalidAccessToken
	}

	validationErr := validations.Validate(in, v.schemas)
	if validationErr != nil {
		return nil, nil, ErrValidationFailed
	}

	return v.next.Image(in, ctx)
}

// Upload implements models.MioService.
func (v *validationService) Upload(in *models.UploadInput, ctx *gin.Context) *vendora.BaseResult {
	tag := models.ObjectTag(in.Tag)
	in.Tag = string(tag)

	file, err := ctx.FormFile("file")
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{ErrBadUploadRequest.Error()},
			Status: http.StatusBadRequest,
		}
	}

	ct, ok := file.Header["Content-Type"]
	if !ok {
		return &vendora.BaseResult{
			Errors: []string{ErrBadUploadRequest.Error()},
			Status: http.StatusBadRequest,
		}
	}

	in.ContentType = ct[0]
	in.Size = file.Size

	validationError := validations.Validate(in, v.schemas)
	if validationError != nil {
		return validationError
	}

	return v.next.Upload(in, ctx)
}

func NewValidationService(opts config.StorageOptions, srv models.MioService, secret, path string) (models.MioService, error) {
	schema := make(map[string][]byte)
	err := validations.LoadSchema(path, schema)
	if err != nil {
		return nil, err
	}

	return &validationService{
		schemas: schema,
		secret:  secret,
		opts:    opts,
		next:    srv,
	}, nil
}
