package user

import (
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
