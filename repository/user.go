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

// FindByUserName is repository for find a user with username
func (r *userRepo) FindByUserName(username string) (*models.UserEntity, error) {
	user := &models.UserEntity{}

	if err := r.Model().Where("user_name = ?", username).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID implements models.UserRepository.
func (r *userRepo) FindByID(id uint) (*models.UserEntity, error) {
	user := &models.UserEntity{}

	err := r.Model().Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// EditUser implements models.UserRepository.
func (r *userRepo) EditUser(id uint, in *models.EditUserInput) (*models.UserEntity, error) {
	user := &models.UserEntity{}

	err := r.Model().Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	user.FirstName = in.FirstName
	user.LastName = in.LastName
	user.Mobile = in.Mobile

	err = r.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUserWithRole implements models.UserRepository.
func (r *userRepo) GetAllUserWithRole(role string) ([]models.UserEntity, error) {
	var users []models.UserEntity

	if strings.ToLower(role) == "all" {
		err := r.Model().Find(users).Error
		if err != nil {
			return nil, err
		}
	} else if strings.ToLower(role) == "user" {
		// err := r.Model().Where(query interface{}, args ...interface{})
	}

	return users, nil
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
