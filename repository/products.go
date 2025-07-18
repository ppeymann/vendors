package repository

import (
	"github.com/ppeymann/vendors.git/models"
	"gorm.io/gorm"
)

type productsRepo struct {
	pg       *gorm.DB
	database string
	table    string
}

func NewProductsRepo(db *gorm.DB, database string) models.ProductRepository {
	return &productsRepo{
		pg:       db,
		database: database,
		table:    "products_entities",
	}
}

// Create implements models.ProductRepository.
func (r *productsRepo) Create(in *models.ProductInput, userID uint) (*models.ProductEntity, error) {
	// first we find user
	userRepo := NewUserRepo(r.pg, r.database)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// create product entity
	pr := &models.ProductEntity{
		Model:            gorm.Model{},
		Title:            in.Title,
		UserID:           userID,
		Description:      in.Description,
		Slug:             "",
		ShortDescription: in.ShortDescription,
		CategoryID:       in.CategoryID,
		Price:            in.Price,
		DiscountPrice:    in.DiscountPrice,
		Stock:            in.Stock,
		SKU:              in.SKU,
		Images:           in.Images,
		Tags:             in.Tags,
		Rating:           0,
		Active:           "DR",
	}

}
