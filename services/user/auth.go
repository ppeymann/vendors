package user

import "github.com/ppeymann/vendors.git/models"

type authService struct {
	next models.UserService
}

func NewAuthService(srv models.UserService) models.UserService {
	return &authService{
		next: srv,
	}
}
