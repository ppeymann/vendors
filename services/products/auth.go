package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
	"github.com/thoas/go-funk"
)

type authService struct {
	next models.ProductService
}

func (a *authService) GetProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	return a.next.GetProduct(ctx, id)
}

// Add implements models.ProductService.
func (a *authService) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	claims, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{vendora.ErrUnAuthorization.Error()},
			Status: http.StatusOK,
		}
	}

	if !funk.Contains(claims.Roles, "SELLER") {
		return &vendora.BaseResult{
			Errors: []string{models.ErrPermissionDenied.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.Add(ctx, in)
}

func NewAuthService(srv models.ProductService) models.ProductService {
	return &authService{
		next: srv,
	}
}
