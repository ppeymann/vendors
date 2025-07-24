package models

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	vendora "github.com/ppeymann/vendors.git"
	"gorm.io/gorm"
)

type (
	// ActiveStatus is activate status for user
	ActiveStatus string

	// ProductService represents method signatures for api Product endpoint.
	// so any object that stratifying this interface can be used as Product service for api endpoint.
	ProductService interface {
		// Add New product
		Add(ctx *gin.Context, in *ProductInput) *vendora.BaseResult

		// GetProduct with specific ID
		GetProduct(ctx *gin.Context, id uint) *vendora.BaseResult

		// GetByTags get all products by tags
		GetByTags(ctx *gin.Context, tags []string) *vendora.BaseResult

		// EditProduct is service for edit a product with specific ID
		EditProduct(ctx *gin.Context, id uint, in *ProductInput) *vendora.BaseResult

		// DeleteProduct with specific ID
		DeleteProduct(ctx *gin.Context, id uint) *vendora.BaseResult
	}

	// ProductRepository represents method signatures for Product domain repository.
	// so any object that stratifying this interface can be used as Product domain repository.
	ProductRepository interface {
		// Create a new Product
		Create(in *ProductInput, userID uint) (*ProductEntity, error)

		// FindByID with specific ID and user ID
		FindByID(id uint) (*ProductEntity, error)

		// FindByTags is repo for find all product that have same tags
		FindByTags(tag []string) ([]*ProductEntity, error)

		// UpdateProduct product
		UpdateProduct(in *ProductInput, id, userID uint) (*ProductEntity, error)

		// Update product
		Update(pr *ProductEntity) error

		// DeleteProduct with specific
		DeleteProduct(id, userID uint) error

		// BaseRepository .
		vendora.BaseRepository
	}

	// ProductHandler represents method signatures for Product handlers.
	// so any object that stratifying this interface can be used as Product handlers.
	ProductHandler interface {
		// Add Handler
		Add(ctx *gin.Context)

		// GetProduct Handler
		GetProduct(ctx *gin.Context)

		// GetByTags handler
		GetByTags(ctx *gin.Context)

		// EditProduct handler
		EditProduct(ctx *gin.Context)

		// DeleteProduct .
		DeleteProduct(ctx *gin.Context)
	}

	// ProductEntity is model for product
	//
	// swagger: model ProductEntity
	ProductEntity struct {
		gorm.Model

		// Title
		Title string `json:"title" gorm:"column:title"`

		// UserID
		UserID uint `json:"user_id" gorm:"column:user_id"`

		// Description
		Description string `json:"desc" gorm:"column:description"`

		// Slug
		Slug string `json:"slug" gorm:"column:slug"`

		// ShortDescription
		ShortDescription string `json:"short_desc" gorm:"column:short_desc"`

		// CategoryID
		CategoryID int64 `json:"category_id" gorm:"column:category_id"`

		// Price is total price
		Price float64 `json:"price" gorm:"column:price"`

		// DiscountPrice
		DiscountPrice float64 `json:"discount_price" gorm:"column:discount_price"`

		// Stock
		Stock int64 `json:"stock" gorm:"column:stock"`

		// SKU is Stock keeping unit
		SKU string `json:"sku" gorm:"column:sku"`

		// Images
		Images pq.StringArray `json:"images" gorm:"column:images;type:text[]"`

		// Tags
		Tags pq.StringArray `json:"tags" gorm:"column:tags;type:text[]"`

		// Rating
		Rating uint32 `json:"rating" gorm:"column:rating"`

		// Active is change from Admin
		Active ActiveStatus `json:"active" gorm:"column:active;default:false"`
	}

	// ProductInput
	//
	// swagger: model ProductInput
	ProductInput struct {
		Title            string         `json:"title"`
		Description      string         `json:"description"`
		ShortDescription string         `json:"short_desc"`
		CategoryID       int64          `json:"category_id"`
		Price            float64        `json:"price"`
		DiscountPrice    float64        `json:"discount_price"`
		Stock            int64          `json:"stock"`
		SKU              string         `json:"sku"`
		Images           pq.StringArray `json:"images" gorm:"images;type:text[]"`
		Tags             pq.StringArray `json:"tags" gorm:"tags;type:text[]"`
	}

	// TagsInput tags input
	TagsInput struct {
		Tags pq.StringArray `json:"tags" gorm:"tags;type:text[]"`
	}
)

const (
	Draft    ActiveStatus = "DR"
	Suspend  ActiveStatus = "SUS"
	Activate ActiveStatus = "AC"
)
