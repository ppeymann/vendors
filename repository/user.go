package repository

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ppeymann/vendors.git/models"
	"gorm.io/gorm"
)

type userRepo struct {
	pg       *gorm.DB
	database string
	table    string
}

func NewUserRepo(db *gorm.DB, database string) models.UserRepository {
	return &userRepo{
		pg:       db,
		database: database,
		table:    "user_entities",
	}
}

func (r *userRepo) Create(in *models.AuthInput) (*models.UserEntity, error) {
	user := &models.UserEntity{
		Model:    gorm.Model{},
		UserName: in.UserName,
		Password: in.Password,
		Roles:    []string{"USER"},
	}

	// create Account
	err := r.pg.Transaction(func(tx *gorm.DB) error {
		if res := r.Model().Create(user).Error; res != nil {
			str := res.(*pgconn.PgError).Message
			if strings.Contains(str, "duplicate key value") {
				return models.ErrAccountExist
			}
			return res
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Update(user *models.UserEntity) error {
	return r.pg.Save(user).Error
}

func (r *userRepo) Migrate() error {
	err := r.pg.AutoMigrate(&models.RefreshTokenEntity{})
	if err != nil {
		return err
	}

	return r.pg.AutoMigrate(&models.UserEntity{})
}

func (r *userRepo) Model() *gorm.DB {
	return r.pg.Model(&models.UserEntity{})
}

func (r *userRepo) Name() string {
	return r.table
}
