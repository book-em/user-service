package api

import (
	"bookem-user-service/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	FindByUsernameOrEmail(username, email string) *domain.User // TODO: Make this return (user, error)
	FindById(id int) (*domain.User, error)
	Update(user *domain.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByUsernameOrEmail(username, email string) *domain.User {
	var user domain.User
	err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *repository) FindById(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
