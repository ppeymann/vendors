package grpcservices

import (
	"context"

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
	in := &models.AuthInput{
		UserName: req.GetUserName(),
		Password: req.GetPassword(),
	}

	user, err := s.repo.Create(in)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateResponse{
		Id:        uint64(user.ID),
		UserName:  user.UserName,
		Password:  user.Password,
		Mobile:    user.Mobile,
		Suspended: user.Suspended,
		Roles:     user.Roles,
	}, nil
}
