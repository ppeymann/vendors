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

		// Login method is for log in site
		Login(ctx *gin.Context, in *AuthInput) *vendora.BaseResult

		// User is for get user information with specific token
		User(ctx *gin.Context) *vendora.BaseResult

		// EditUser is for edit user entity (mobile, first name, ...)
		EditUser(ctx *gin.Context, in *EditUserInput) *vendora.BaseResult

		// ActiveDeActiveSuspended .
		ActiveDeActiveSuspended(ctx *gin.Context) *vendora.BaseResult

		// GetAllUserWithRole this method is for admin to get all user with specific role
		GetAllUserWithRole(ctx *gin.Context, role string) *vendora.BaseResult
	}

	// UserRepository represents method signatures for user domain repository.
	// so any object that stratifying this interface can be used as user domain repository.
	UserRepository interface {
		// Create is for create a new user
		Create(in *AuthInput) (*UserEntity, error)

		// Update is for updating user information
		Update(user *UserEntity) error

		// FindByUserName
		FindByUserName(username string) (*UserEntity, error)

		// FindByID
		FindByID(id uint) (*UserEntity, error)

		// EditUser
		EditUser(id uint, in *EditUserInput) (*UserEntity, error)

		// GetAllUserWithRole
		GetAllUserWithRole(role string) ([]UserEntity, error)

		// BaseRepository (migrate, models,...)
		vendora.BaseRepository
	}

	// UserHandler represents method signatures for user handlers.
	// so any object that stratifying this interface can be used as user handlers.
	UserHandler interface {
		// Register is handler for sign up
		Register(ctx *gin.Context)

		// Login is handler for log in
		Login(ctx *gin.Context)

		// User is handler for get user information
		User(ctx *gin.Context)

		// EditUser
		EditUser(ctx *gin.Context)

		// GetAllUserWithRole
		GetAllUserWithRole(ctx *gin.Context)

		// ActiveDeActiveSuspended
		ActiveDeActiveSuspended(ctx *gin.Context)
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
		Mobile string `json:"mobile" gorm:"column:mobile;index"`

		// Suspended uses as determination flag for account suspension situation
		Suspended bool `json:"suspended" gorm:"column:suspended;index"`

		// Roles contains account access level permissions
		Roles pq.StringArray `json:"roles" gorm:"type:varchar(64)[]"`

		// Password
		Password string `json:"-" gorm:"column:password;not null;size:100" mapstructure:"password"`

		// Tokens list of current user active session
		Tokens []RefreshTokenEntity `json:"-" gorm:"foreignKey:AccountID;references:ID"`

		// FirstName
		FirstName string `json:"first_name" gorm:"column:first_name" mapstructure:"first_name"`

		// LastName
		LastName string `json:"last_name" gorm:"column:last_name" mapstructure:"last_name"`

		// Balance
		Balance float64 `json:"balance" gorm:"column:balance;default:0.00"`

		// Products just not null for seller
		Products []ProductEntity `json:"products" gorm:"foreignKey:UserID;references:ID"`
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

	// EditUserInput is input for editing user entity
	EditUserInput struct {
		FirstName string `json:"first_name" mapstructure:"first_name"`
		LastName  string `json:"last_name" mapstructure:"last_name"`
		Mobile    string `json:"mobile" mapstructure:"mobile"`
	}
)
