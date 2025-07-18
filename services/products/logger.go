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

func NewLoggerService(logger log.Logger, srv models.ProductService) models.ProductService {
	return &loggerService{
		logger: logger,
		next:   srv,
	}
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
