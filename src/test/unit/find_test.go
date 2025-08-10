package test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestFindById_Success(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)

	user := defaultUser
	user.ID = id

	mockRepo.On("FindById", id).Return(user, nil)
	userGot, err := svc.FindById(id)

	assert.NoError(t, err)
	assert.NotNil(t, userGot)
	assert.Equal(t, user.ID, userGot.ID)

}

func TestFindById_UserNotFound(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	id := uint(1)

	user := defaultUser
	user.ID = id

	mockRepo.On("FindById", id).Return(nil, fmt.Errorf("no such user"))
	userGot, err := svc.FindById(id)

	assert.Error(t, err)
	assert.Nil(t, userGot)
}
