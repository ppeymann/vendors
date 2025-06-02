package grpcservices

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/ppeymann/vendors.git/models"
	userpb "github.com/ppeymann/vendors.git/proto/user"
)

type userService struct {
	userpb.UnimplementedUserServiceServer
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) userpb.UserServiceServer {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	in := &models.AuthInput{}

	_ = mapstructure.Decode(req, in)

	return nil, nil

}
