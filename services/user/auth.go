package user

import (
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type authService struct {
	next models.UserService
}

func NewAuthService(srv models.UserService) models.UserService {
	return &authService{
		next: srv,
	}
}

// Register implements models.UserService.
func (a *authService) Register(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	return a.next.Register(ctx, in)
}
