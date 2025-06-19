package mio

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/auth"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/utils"
)

type (
	authService struct {
		opts config.StorageOptions
		next models.MioService
	}
)

func NewAuthService(opts config.StorageOptions, srv models.MioService) models.MioService {
	return &authService{
		opts: opts,
		next: srv,
	}
}

// Download implements models.MioService.
func (a *authService) Download(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	return a.next.Download(in, ctx)
}

// Image implements models.MioService.
func (a *authService) Image(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	return a.next.Image(in, ctx)
}

// Upload implements models.MioService.
func (a *authService) Upload(in *models.UploadInput, ctx *gin.Context) *vendora.BaseResult {
	claims := &auth.Claims{}
	err := utils.CatchClaims(ctx, claims)
	if err != nil {
		return &vendora.BaseResult{
			Status: http.StatusOK,
			Errors: []string{ErrInvalidAccessToken.Error()},
		}
	}

	in.Claims = claims

	return a.next.Upload(in, ctx)
}
