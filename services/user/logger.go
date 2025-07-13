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

// Login implements models.UserService.
func (l *loggerService) Login(ctx *gin.Context, in *models.AuthInput) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Login",
			"errors", strings.Join(result.Errors, ","),
			"user_name", in.UserName,
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.Login(ctx, in)
}

// User implements models.UserService.
func (l *loggerService) User(ctx *gin.Context) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "User",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.User(ctx)
}

// EditUser implements models.UserService.
func (l *loggerService) EditUser(ctx *gin.Context, in *models.EditUserInput) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "EditUser",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.EditUser(ctx, in)
}

// GetAllUserWithRole implements models.UserService.
func (l *loggerService) GetAllUserWithRole(ctx *gin.Context, role string) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetAllUserWithRole",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"role", role,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.GetAllUserWithRole(ctx, role)
}

// ActiveDeActiveSuspended implements models.UserService.
func (l *loggerService) ActiveDeActiveSuspended(ctx *gin.Context) (result *vendora.BaseResult) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "ActiveDeActiveSuspended",
			"errors", strings.Join(result.Errors, " ,"),
			"result", result,
			"client_ip", ctx.ClientIP(),
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.next.ActiveDeActiveSuspended(ctx)
}
