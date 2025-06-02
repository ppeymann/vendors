package user

import (
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
	userpb "github.com/ppeymann/vendors.git/proto/user"
)

type service struct {
	client userpb.UserServiceClient
	conf   *config.Configuration
}

func NewService(client userpb.UserServiceClient, conf *config.Configuration) models.UserService {
	return &service{
		client: client,
		conf:   conf,
	}
}

// TODO: Add this method
func (s *service) Register(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	req := &userpb.CreateRequest{
		UserName: in.UserName,
		Password: in.Password,
	}

	_, err := s.client.Create(ctx, req)
	if err != nil {
		return nil
	}

	return nil
}
