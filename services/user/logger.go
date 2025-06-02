package user

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type loggerService struct {
	next   models.UserService
	logger log.Logger
}

func NewLoggerService(srv models.UserService, logger log.Logger) models.UserService {
	return &loggerService{
		next:   srv,
		logger: logger,
	}
}

func (l *loggerService) Register(ctx *gin.Context, in *models.AuthInput) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Register",
			"errors", strings.Join(result.Errors, ","),
			"user_name", in.UserName,
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.Register(ctx, in)
}
