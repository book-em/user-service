package test

import (
	"context"
	"fmt"
	"testing"

	"bookem-user-service/client/roomclient"
	domain "bookem-user-service/domain"

	assert "github.com/stretchr/testify/assert"
)

func TestDelete_Success(t *testing.T) {
	svc, mockRepo, mockRoomClient := createTestService()

	id := uint(1)
	callerID := uint(1)
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	mockRoomClient.On("GetPendingGuestReservations", context.Background(), user).Return([]roomclient.ReservationDTO{}, nil)

	err := svc.Delete(context.Background(), callerID, id)

	assert.NoError(t, err)
}

func TestDelete_TriedToDeleteSomeoneElse(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)
	callerID := uint(2)
	user := defaultUser
	user.ID = id
	mockRepo.On("FindById", id).Return(user, nil)

	err := svc.Delete(context.Background(), callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)
	callerID := uint(1)
	mockRepo.On("FindById", id).Return(nil, fmt.Errorf("user not found"))

	err := svc.Delete(context.Background(), callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestDelete_GuestHasPendingReservations(t *testing.T) {
	svc, mockRepo, mockRoomClient := createTestService()

	id := uint(1)
	callerID := uint(1)
	user := defaultUser
	user.ID = id
	user.Role = domain.Guest
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := roomclient.ReservationDTO{}
	mockRoomClient.On("GetPendingGuestReservations", context.Background(), user).Return([]roomclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrGuestHasReservations, err)
}

func TestDelete_HostHasPendingReservations(t *testing.T) {
	svc, mockRepo, mockRoomClient := createTestService()

	id := uint(1)
	callerID := uint(1)
	user := defaultUser
	user.ID = id
	user.Role = domain.Host
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()
	reservation := roomclient.ReservationDTO{}
	mockRoomClient.On("GetActiveHostReservations", context.Background(), user).Return([]roomclient.ReservationDTO{reservation}, nil)

	err := svc.Delete(context.Background(), callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrHostHasReservations, err)
}

func TestDelete_TriedDeletingAdmin(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)
	callerID := uint(1)
	user := defaultUser
	user.ID = id
	user.Role = domain.Admin
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()

	err := svc.Delete(context.Background(), callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCannotDeleteAdmin, err)
}
