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

// Add implements models.ProductService.
func (v *validationService) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	err := validations.Validate(in, v.schema)
	if err != nil {
		return err
	}

	return v.next.Add(ctx, in)
}
