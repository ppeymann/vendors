package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type service struct {
	repo models.ProductRepository
}

func NewService(repo models.ProductRepository) models.ProductService {
	return &service{
		repo: repo,
	}
}

// Add implements models.ProductService.
func (s *service) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	claims, _ := vendora.CheckAuth(ctx)

	pr, err := s.repo.Create(in, claims.Subject)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: pr,
	}
}
