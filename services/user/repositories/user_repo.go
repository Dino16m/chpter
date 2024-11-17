package repositories

import (
	"github.com/dino16m/chpter/user/apperrs"
	"github.com/dino16m/chpter/user/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {

	err := r.db.Create(user).Error

	if err != nil {
		return apperrs.NewServerError("An error occurred while creating a new user", err)
	}
	return nil
}

func (r *UserRepository) FindById(id uint) (models.User, error) {
	var user models.User
	tx := r.db.Model(&models.User{}).First(&user, id)
	return apperrs.WrapNotFound(user, tx.Error, "User not found")
}
