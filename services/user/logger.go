package user

import (
	"github.com/go-kit/log"
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
