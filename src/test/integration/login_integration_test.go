package test

import (
	"bookem-user-service/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func loginUser(username_or_email string, password string) (*http.Response, error) {
	dto := domain.LoginDTO{
		UsernameOrEmail: username_or_email,
		Password:        password,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL+"login", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func loginUser2(username_or_email string, password string) string {
	resp, _ := loginUser(username_or_email, password)
	body := resp.Body.Close().Error()
	var token domain.JWTDTO
	json.Unmarshal([]byte(body), &token)
	return token.Jwt
}

func TestIntegration_LoginUser(t *testing.T) {
	registerUser("user_03", "1234", ROLE_GUEST)
	resp, err := loginUser("user_03", "1234")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	registerUser("user_04@gmail.com", "1234", ROLE_GUEST)
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
	registerUser("user_05", "password1234", ROLE_GUEST)
	resp, err := loginUser("user_05", "a")

	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
