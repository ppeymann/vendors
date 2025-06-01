package models

import (
	"errors"
	"time"

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
		// Register(ctx *gin.Context, in *AuthInput) *vendora.BaseResult
	}

	// UserRepository represents method signatures for user domain repository.
	// so any object that stratifying this interface can be used as user domain repository.
	UserRepository interface {
	}

	// UserHandler represents method signatures for user handlers.
	// so any object that stratifying this interface can be used as user handlers.
	UserHandler interface {
	}

	AuthInput struct {
		UserName string `json:"user_name" mapstructure:"user_name"`
		Password string `json:"password" mapstructure:"password"`
	}

	UserEntity struct {
		gorm.Model `swaggerignore:"true"`

		// UserName
		UserName string `json:"user_name" gorm:"column:user_name;index;unique" mapstructure:"user_name"`

		// Password
		Password string `json:"-" gorm:"column:password;not null;size:100" mapstructure:"password"`

		// FullName
		Fullname string `json:"full_name" gorm:"column:full_name" mapstructure:"full_name"`
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
