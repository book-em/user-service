package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_LoginUser(t *testing.T) {
	registerUser("user_03", "1234", domain.Guest)
	resp, err := loginUser("user_03", "1234")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	registerUser("user_04@gmail.com", "1234", domain.Guest)
	resp, err = loginUser("user_04@gmail.com", "1234")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegration_LoginBadUser(t *testing.T) {
	resp, err := loginUser("", "1234")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = loginUser("user_03@failmail.org", "1234")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_LoginBadPassword(t *testing.T) {
	registerUser("user_05", "password1234", domain.Guest)
	resp, err := loginUser("user_05", "a")

	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
