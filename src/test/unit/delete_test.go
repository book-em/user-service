package test

import (
	"bookem-user-service/client/reservationclient"
	"bookem-user-service/client/roomclient"
	domain "bookem-user-service/domain"
	"context"
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestDelete_Success(t *testing.T) {
	svc, mockRepo, _, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest

	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	mockReservationClient.On("GetActiveGuestReservations", jwt).Return([]roomclient.ReservationDTO{}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.NoError(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	svc, mockRepo, _, _ := createTestService()

	id := uint(1)
	jwt := "token"
	mockRepo.On("FindById", id).Return(nil, fmt.Errorf("user not found"))

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestDelete_GuestHasActiveReservations(t *testing.T) {
	svc, mockRepo, _, mockReservationClient := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := reservationclient.ReservationDTO{}
	mockReservationClient.On("GetActiveGuestReservations", jwt).Return([]reservationclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrGuestHasReservations, err)
}

func TestDelete_HostHasActiveReservations(t *testing.T) {
	svc, mockRepo, mockRoomClient, _ := createTestService()

	id := uint(1)
	jwt := "token"
	user := defaultUser
	user.ID = id
	user.Role = domain.Host
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := roomclient.ReservationDTO{}
	mockRoomClient.On("GetActiveHostReservations", jwt).Return([]roomclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), id, jwt)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrHostHasReservations, err)
}

func TestDelete_TriedDeletingAdmin(t *testing.T) {
	svc, mockRepo, _, _ := createTestService()

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
