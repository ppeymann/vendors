package models

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	vendora "github.com/ppeymann/vendors.git"
	"gorm.io/gorm"
)

type (
	// ProductService represents method signatures for api Product endpoint.
	// so any object that stratifying this interface can be used as Product service for api endpoint.
	ProductService interface {
		Add(ctx *gin.Context, in *ProductEntity) *vendora.BaseResult
	}

	// ProductRepository represents method signatures for Product domain repository.
	// so any object that stratifying this interface can be used as Product domain repository.
	ProductRepository interface {
		Create()
	}

	// ProductHandler represents method signatures for Product handlers.
	// so any object that stratifying this interface can be used as Product handlers.
	ProductHandler interface{}

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
		CategoryID string `json:"category_id" gorm:"column:category_id"`

		// Price is total price
		Price int64 `json:"price" gorm:"column:price"`

		// DiscountPrice
		DiscountPrice int64 `json:"discount_price" gorm:"column:discount_price"`

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
	}
)
