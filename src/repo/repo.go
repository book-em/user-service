package api

import (
	"bookem-user-service/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	FindByUsernameOrEmail(username, email string) (*domain.User, error)
	FindById(id uint) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint)
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

func (r *repository) FindByUsernameOrEmail(username, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindById(id uint) (*domain.User, error) {
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

func (r *repository) Delete(id uint) {
	r.db.Delete(&domain.User{}, id)
}
