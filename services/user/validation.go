package user

import (
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
	validations "github.com/ppeymann/vendors.git/validation"
)

type validationService struct {
	next   models.UserService
	schema map[string][]byte
}

func NewValidationService(srv models.UserService, path string) (models.UserService, error) {
	schema := make(map[string][]byte)

	// Load the schema from the specified path
	err := validations.LoadSchema(path, schema)
	if err != nil {
		return nil, err
	}

	return &validationService{
		next:   srv,
		schema: schema,
	}, nil
}

func (v *validationService) Register(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	err := validations.Validate(in, v.schema)
	if err != nil {
		return err
	}

	return v.next.Register(ctx, in)
}

// Login implements models.UserService.
func (v *validationService) Login(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	err := validations.Validate(in, v.schema)
	if err != nil {
		return err
	}

	return v.next.Register(ctx, in)
}
