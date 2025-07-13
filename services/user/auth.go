package user

import (
	"net/http"

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

// Login implements models.UserService.
func (a *authService) Login(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	return a.next.Login(ctx, in)
}

// User implements models.UserService.
func (a *authService) User(ctx *gin.Context) *vendora.BaseResult {
	_, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.User(ctx)
}

// EditUser implements models.UserService.
func (a *authService) EditUser(ctx *gin.Context, in *models.EditUserInput) *vendora.BaseResult {
	_, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.EditUser(ctx, in)
}

// GetAllUserWithRole implements models.UserService.
func (a *authService) GetAllUserWithRole(ctx *gin.Context, role string) *vendora.BaseResult {
	_, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.GetAllUserWithRole(ctx, role)
}

// ActiveDeActiveSuspended implements models.UserService.
func (a *authService) ActiveDeActiveSuspended(ctx *gin.Context) *vendora.BaseResult {
	_, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.ActiveDeActiveSuspended(ctx)
}
