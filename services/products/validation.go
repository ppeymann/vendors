package products

import (
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
	validations "github.com/ppeymann/vendors.git/validation"
)

type validationService struct {
	schema map[string][]byte
	next   models.ProductService
}

func (v *validationService) DeleteProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	return v.next.DeleteProduct(ctx, id)
}

func (v *validationService) EditProduct(ctx *gin.Context, id uint, in *models.ProductInput) *vendora.BaseResult {
	err := validations.Validate(in, v.schema)
	if err != nil {
		return err
	}

	return v.next.EditProduct(ctx, id, in)
}

func (v *validationService) GetByTags(ctx *gin.Context, tags []string) *vendora.BaseResult {
	return v.next.GetByTags(ctx, tags)
}

func (v *validationService) GetProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	return v.next.GetProduct(ctx, id)
}

// Add implements models.ProductService.
func (v *validationService) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	err := validations.Validate(in, v.schema)
	if err != nil {
		return err
	}

	return v.next.Add(ctx, in)
}

func NewValidationsService(path string, srv models.ProductService) (models.ProductService, error) {
	schemas := make(map[string][]byte)

	err := validations.LoadSchema(path, schemas)
	if err != nil {
		return nil, err
	}

	return &validationService{
		schema: schemas,
		next:   srv,
	}, nil
}
