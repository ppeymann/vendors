package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	vendora "github.com/ppeymann/vendors.git"
	"gorm.io/gorm"
)

// Errors
var (
	ErrAccountExist     error = errors.New("account with specified params already exists")
	ErrSignInFailed     error = errors.New("account not found or password error")
	ErrPermissionDenied error = errors.New("specified role is not available for user")
	ErrAccountNotExist  error = errors.New("specified account does not exist")
)

type (
	// UserService represents method signatures for api user endpoint.
	// so any object that stratifying this interface can be used as user service for api endpoint.
	UserService interface {
		// Register is for sign up new user
		Register(ctx *gin.Context, in *AuthInput) *vendora.BaseResult
	}

	// UserRepository represents method signatures for user domain repository.
	// so any object that stratifying this interface can be used as user domain repository.
	UserRepository interface {
		// Create is for create a new user
		Create(in *AuthInput) (*UserEntity, error)

		// Update is for updating user information
		Update(user *UserEntity) error

		// BaseRepository (migrate, models,...)
		vendora.BaseRepository
	}

	// UserHandler represents method signatures for user handlers.
	// so any object that stratifying this interface can be used as user handlers.
	UserHandler interface {
		// Register is handler for sign up
		Register(ctx *gin.Context)
	}

	AuthInput struct {
		UserName string `json:"user_name" mapstructure:"user_name"`
		Password string `json:"password" mapstructure:"password"`
	}

	UserEntity struct {
		gorm.Model `swaggerignore:"true"`

		// UserName
		UserName string `json:"user_name" gorm:"column:user_name;index;unique" mapstructure:"user_name"`

		// Mobile phone number of account owner
		Mobile string `json:"mobile" gorm:"column:mobile;index;unique"`

		// Suspended uses as determination flag for account suspension situation
		Suspended bool `json:"suspended" gorm:"column:suspended;index"`

		// Roles contains account access level permissions
		Roles pq.StringArray `json:"roles" gorm:"type:varchar(64)[]"`

		// Password
		Password string `json:"-" gorm:"column:password;not null;size:100" mapstructure:"password"`

		// Tokens list of current user active session
		Tokens []RefreshTokenEntity `json:"-" gorm:"foreignKey:AccountID;references:ID"`
	}

	TokenBundlerOutput struct {
		// Token is string that hashed by paseto
		Token string `json:"token"`

		// Refresh is string that for refresh old token
		Refresh string `json:"refresh"`

		// Expire is time for expire token
		Expire time.Time `json:"expire"`
	}

	// RefreshTokenEntity is entity to store accounts active session
	RefreshTokenEntity struct {
		gorm.Model
		AccountID uint
		TokenId   string `json:"token_id" gorm:"column:token_id;index"`
		UserAgent string `json:"user_agent" gorm:"column:user_agent"`
		IssuedAt  int64  `json:"issued_at" bson:"issued_at" gorm:"column:issued_at"`
		ExpiredAt int64  `json:"expired_at" bson:"expired_at" gorm:"column:expired_at"`
	}
)
