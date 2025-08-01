package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_RegisterUser(t *testing.T) {
	resp, err := registerUser("user1", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestIntegration_RegisterUserDuplicate(t *testing.T) {
	resp, err := registerUser("user_02", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = registerUser("user_02", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	resp, err = registerUser("user_02@gmail.com", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}
