package repository

import (
	"errors"
	"github.com/ppeymann/vendors.git/models"
	"gorm.io/gorm"
)

type productsRepo struct {
	pg       *gorm.DB
	database string
	table    string
}

func (r *productsRepo) Update(pr *models.ProductEntity) error {
	return r.Model().Save(pr).Error
}

func (r *productsRepo) UpdateProduct(in *models.ProductInput, id, userID uint) (*models.ProductEntity, error) {
	pr, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}

	if pr.UserID != userID {
		return nil, errors.New("permission denied")
	}

	pr.Title = in.Title
	pr.Description = in.Description
	pr.Stock = in.Stock
	pr.Price = in.Price
	pr.Tags = in.Tags
	pr.SKU = in.SKU
	pr.DiscountPrice = in.DiscountPrice
	pr.Images = in.Images
	pr.ShortDescription = in.ShortDescription
	pr.CategoryID = in.CategoryID

	err = r.Update(pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *productsRepo) FindByTags(tag []string) ([]*models.ProductEntity, error) {
	var products []*models.ProductEntity

	if err := r.Model().Where("tag IN (?)", tag).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productsRepo) FindByID(id uint) (*models.ProductEntity, error) {
	pr := &models.ProductEntity{}
	err := r.Model().Where("id = ?", id).First(pr).Error
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// Migrate implements models.ProductRepository.
func (r *productsRepo) Migrate() error {
	return r.pg.AutoMigrate(&models.ProductEntity{})
}

// Model implements models.ProductRepository.
func (r *productsRepo) Model() *gorm.DB {
	return r.pg.Model(&models.ProductEntity{})
}

// Name implements models.ProductRepository.
func (r *productsRepo) Name() string {
	return r.table
}

// Create implements models.ProductRepository.
func (r *productsRepo) Create(in *models.ProductInput, userID uint) (*models.ProductEntity, error) {
	// first we find user
	userRepo := NewUserRepo(r.pg, r.database)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Suspended {
		return nil, errors.New("user is not active")
	}

	// create product entity
	pr := &models.ProductEntity{
		Model:            gorm.Model{},
		Title:            in.Title,
		UserID:           user.ID,
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

	err = r.Model().Create(pr).Error
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func NewProductsRepo(db *gorm.DB, database string) models.ProductRepository {
	return &productsRepo{
		pg:       db,
		database: database,
		table:    "products_entities",
	}
}
