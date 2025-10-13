package test

import (
	"bookem-user-service/client/reservationclient"
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	"context"

	mock "github.com/stretchr/testify/mock"
)

func createTestService() (service.Service, *MockRepo, *MockRoomClient, *MockReservationClient) {
	mockRepo := new(MockRepo)
	mockRoomClient := new(MockRoomClient)
	mockReservationClient := new(MockReservationClient)

	svc := service.NewService(mockRepo, mockRoomClient, mockReservationClient)

	return svc, mockRepo, mockRoomClient, mockReservationClient
}

// ---------------------------------------------- Mock repo

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindByUsernameOrEmail(username, email string) (*domain.User, error) {
	args := m.Called(username, email)
	user, _ := args.Get(0).(*domain.User)
	err := args.Error(1)
	return user, err
}

func (m *MockRepo) FindByUsernameOrEmailNotId(username, email string, id uint) (*domain.User, error) {
	args := m.Called(username, email, id)
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

// ---------------------------------------------- Mock room client

type MockRoomClient struct {
	mock.Mock
}

// ---------------------------------------------- Mock reservation client

type MockReservationClient struct {
	mock.Mock
}

func (m *MockReservationClient) GetActiveGuestReservations(ctx context.Context, jwt string) ([]reservationclient.ReservationDTO, error) {
	args := m.Called(ctx, jwt)
	reservations, _ := args.Get(0).([]reservationclient.ReservationDTO)
	return reservations, args.Error(1)
}

func (m *MockReservationClient) GetActiveHostReservations(ctx context.Context, jwt string) ([]reservationclient.ReservationDTO, error) {
	args := m.Called(ctx, jwt)
	reservations, _ := args.Get(0).([]reservationclient.ReservationDTO)
	return reservations, args.Error(1)
}

// ---------------------------------------------- Mock data

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
