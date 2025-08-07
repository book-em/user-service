package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindById(t *testing.T) {
	resp, _ := registerUser("user_30", "1234", domain.Guest)
	id := getUserFromRegister(resp).Id
	resp, err := findUserById(id)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegration_FindById_NotFound(t *testing.T) {
	resp, _ := registerUser("user_31", "1234", domain.Guest)
	resp, err := findUserById(0)

	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
