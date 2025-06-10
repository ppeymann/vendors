package mio

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type loggerService struct {
	logger log.Logger
	next   models.MioService
}

// Download implements models.MioService.
func (l *loggerService) Download(in *models.DownloadInput, ctx *gin.Context) (result []byte, file *models.StorageEntity, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "download",
			"account_id", in.AccountId,
			"tag", in.Tag,
			"file_id", in.Id,
			"took", time.Since(begin),
			"client_ip", ctx.GetHeader("X-Real-Ip"),
			"err", err,
		)
	}(time.Now())

	return l.next.Download(in, ctx)
}

// Image implements models.MioService.
func (l *loggerService) Image(in *models.DownloadInput, ctx *gin.Context) (result []byte, file *models.StorageEntity, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "image",
			"account_id", in.AccountId,
			"tag", in.Tag,
			"file_id", in.Id,
			"took", time.Since(begin),
			"client_ip", ctx.GetHeader("X-Real-Ip"),
			"err", err,
		)
	}(time.Now())

	return l.next.Image(in, ctx)
}

// Upload implements models.MioService.
func (l *loggerService) Upload(in *models.UploadInput, ctx *gin.Context) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		errs := ""
		if result.Errors != nil {
			errs = strings.Join(result.Errors, " ,")
		}

		_ = l.logger.Log(
			"method", "Upload",
			"user", in.Claims.Subject,
			"tag", in.Tag,
			"size", in.Size,
			"took", time.Since(begin),
			"client_ip", ctx.GetHeader("X-Real-Ip"),
			"err", errs,
		)

	}(time.Now())

	return l.next.Upload(in, ctx)
}

func NewLoggerService(logger log.Logger, srv models.MioService) models.MioService {
	return &loggerService{
		logger: logger,
		next:   srv,
	}
}
