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

func (s *service) DeleteProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	claims, _ := vendora.CheckAuth(ctx)
	err := s.repo.DeleteProduct(id, claims.Subject)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Result: id,
		Status: http.StatusOK,
	}
}

func (s *service) EditProduct(ctx *gin.Context, id uint, in *models.ProductInput) *vendora.BaseResult {
	claims, _ := vendora.CheckAuth(ctx)
	pr, err := s.repo.UpdateProduct(in, id, claims.Subject)
	if err != nil {
		return &vendora.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: pr,
	}
}

func (s *service) GetByTags(ctx *gin.Context, tags []string) *vendora.BaseResult {
	pr, err := s.repo.FindByTags(tags)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status:      http.StatusOK,
		Result:      pr,
		ResultCount: int64(len(pr)),
	}
}

func (s *service) GetProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	pr, err := s.repo.FindByID(id)
	if err != nil {
		return &vendora.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: pr,
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

func NewService(repo models.ProductRepository) models.ProductService {
	return &service{
		repo: repo,
	}
}
