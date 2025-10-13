package test

import (
	"context"
	"fmt"
	"testing"

	domain "bookem-user-service/domain"
	utils "bookem-user-service/util"

	assert "github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	dto := domain.LoginDTO{
		UsernameOrEmail: "user123",
		Password:        "1234",
	}

	pwHashed, _ := utils.HashPassword(dto.Password)

	user := domain.User{
		Username: "user123",
		Password: pwHashed,
		Email:    "user123@gmail.com",
		Name:     "abc",
		Surname:  "def",
		Address:  "Address 123",
		Role:     "guest",
	}

	// Mock

	utils.CreateJWT = func(userID int, username string, role domain.UserRole) (string, error) {
		return "aaa", nil
	}

	mockRepo.On(
		"FindByUsernameOrEmail",
		dto.UsernameOrEmail, dto.UsernameOrEmail,
	).Return(&user, nil)

	// Verify

	jwt, err := svc.Login(context.Background(), dto)

	assert.NoError(t, err)
	assert.NotEqual(t, "", jwt)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := domain.LoginDTO{
		UsernameOrEmail: "user123",
		Password:        "1234",
	}

	mockRepo.On(
		"FindByUsernameOrEmail",
		dto.UsernameOrEmail, dto.UsernameOrEmail,
	).Return(nil, fmt.Errorf("no such user"))

	jwt, err := svc.Login(context.Background(), dto)

	assert.ErrorIs(t, err, domain.ErrLoginFailed)
	assert.Equal(t, "", jwt)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := domain.LoginDTO{
		UsernameOrEmail: "user123",
		Password:        "1234",
	}

	user := domain.User{
		Username: "user123",
		Password: "jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
		Email:    "user123@gmail.com",
		Name:     "abc",
		Surname:  "def",
		Address:  "Address 123",
		Role:     "guest",
	}

	// Note: user.Password should be hashed, so even if it was set to "user123",
	// the password check would fail.

	mockRepo.On(
		"FindByUsernameOrEmail",
		dto.UsernameOrEmail, dto.UsernameOrEmail,
	).Return(&user, nil)

	jwt, err := svc.Login(context.Background(), dto)

	assert.ErrorIs(t, err, domain.ErrLoginFailed)
	assert.Equal(t, "", jwt)
}

func TestLogin_JWTFailed(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	dto := domain.LoginDTO{
		UsernameOrEmail: "user123",
		Password:        "1234",
	}

	pwHashed, _ := utils.HashPassword(dto.Password)

	user := domain.User{
		Username: "user123",
		Password: pwHashed,
		Email:    "user123@gmail.com",
		Name:     "abc",
		Surname:  "def",
		Address:  "Address 123",
		Role:     "guest",
	}

	// Mock

	utils.CreateJWT = func(userID int, username string, role domain.UserRole) (string, error) {
		return "", fmt.Errorf("Some error")
	}

	mockRepo.On(
		"FindByUsernameOrEmail",
		dto.UsernameOrEmail, dto.UsernameOrEmail,
	).Return(&user, nil)

	// Verify

	jwt, err := svc.Login(context.Background(), dto)

	assert.Error(t, err)
	assert.Equal(t, "", jwt)
}
