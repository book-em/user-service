package test

import (
	domain "bookem-user-service/domain"

	mock "github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindByUsernameOrEmail(username, email string) (*domain.User, error) {
	args := m.Called(username, email)
	user, _ := args.Get(0).(*domain.User)
	err := args.Error(1)
	return user, err
}

func (m *MockRepo) FindById(id uint) (*domain.User, error) {
	args := m.Called(uint(id))
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepo) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepo) Delete(id uint) {
	m.Called(id)
}

var defaultUserDTO = &domain.UserCreateDTO{
	Username: "user",
	Password: "pass",
	Email:    "email@mail.com",
	Name:     "name",
	Surname:  "surname",
	Role:     "guest",
	Address:  "Address 123",
}

var defaultUser = &domain.User{
	Username: "user",
	Password: "pass",
	Email:    "email@mail.com",
	Name:     "name",
	Surname:  "surname",
	Role:     "guest",
	Address:  "Address 123",
}
