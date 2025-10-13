package test

import (
	"bookem-user-service/client/reservationclient"
	domain "bookem-user-service/domain"
	"context"
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestDelete_GuestSuccess(t *testing.T) {
	svc, mockRepo, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest

	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	mockReservationClient.On("GetActiveGuestReservations", context.Background(), jwt).Return([]reservationclient.ReservationDTO{}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.NoError(t, err)
}

func TestDelete_HostSuccess(t *testing.T) {
	svc, mockRepo, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Host

	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt).Return([]reservationclient.ReservationDTO{}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.NoError(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)
	jwt := "token"
	mockRepo.On("FindById", id).Return(nil, fmt.Errorf("user not found"))

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestDelete_GuestHasActiveReservations(t *testing.T) {
	svc, mockRepo, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := reservationclient.ReservationDTO{}
	mockReservationClient.On("GetActiveGuestReservations", context.Background(), jwt).Return([]reservationclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrGuestHasReservations, err)
}

func TestDelete_HostHasActiveReservations(t *testing.T) {
	svc, mockRepo, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Host
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := reservationclient.ReservationDTO{}
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt).Return([]reservationclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrHostHasReservations, err)
}

func TestDelete_TriedDeletingAdmin(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Admin
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCannotDeleteAdmin, err)
}
