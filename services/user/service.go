package user

import "github.com/ppeymann/vendors.git/models"

type service struct {
}

func NewService() models.UserService {
	return &service{}
}
