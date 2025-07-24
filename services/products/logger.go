package products

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
	next   models.ProductService
}

func (l *loggerService) EditProduct(ctx *gin.Context, id uint, in *models.ProductInput) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "EditProduct",
			"errors", strings.Join(result.Errors, ","),
			"result", result,
			"id", id,
			"input", in,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.EditProduct(ctx, id, in)
}

func (l *loggerService) GetByTags(ctx *gin.Context, tags []string) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetByTags",
			"errors", strings.Join(result.Errors, ","),
			"result", result,
			"tags", tags,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.GetByTags(ctx, tags)
}

func (l *loggerService) GetProduct(ctx *gin.Context, id uint) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetProduct",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"id", id,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.GetProduct(ctx, id)
}

// Add implements models.ProductService.
func (l *loggerService) Add(ctx *gin.Context, in *models.ProductInput) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Add",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.Add(ctx, in)
}

func NewLoggerService(logger log.Logger, srv models.ProductService) models.ProductService {
	return &loggerService{
		logger: logger,
		next:   srv,
	}
}
