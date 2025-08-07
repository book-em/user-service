package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_PasswordChange_Success(t *testing.T) {
	resp, err := registerUser("user_20", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_20", "1234")

	resp, err = changePassword(jwt, user.Id, "1234", "12345", "12345")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_PasswordChange_WrongUser(t *testing.T) {
	registerUser("user_21_somebody_else", "1234", domain.Guest)
	resp, _ := registerUser("user_21", "1234", domain.Guest)

	user := getUserFromRegister(resp)

	// Somebody else.

	jwt := loginUser2("user_21_somebody_else", "1234")
	resp, err := changePassword(jwt, user.Id, "1234", "12345", "12345")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Signed out user.

	resp, err = changePassword("", user.Id, "1234", "12345", "12345")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestIntegration_PasswordChange_ConfirmPasswordFail(t *testing.T) {
	resp, _ := registerUser("user_23_1", "1234", domain.Guest)

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_23_1", "1234")

	resp, err := changePassword(jwt, user.Id, "1234", "12345", "12345jjjjjjjjj")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_PasswordChange_PasswordNotChanged(t *testing.T) {
	resp, _ := registerUser("user_24_1", "1234", domain.Guest)

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_24_1", "1234")

	resp, err := changePassword(jwt, user.Id, "1234", "1234", "1234")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_PasswordChange_OldPasswordFail(t *testing.T) {
	resp, _ := registerUser("user_25_1", "1234", domain.Guest)

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_25_1", "1234")

	resp, err := changePassword(jwt, user.Id, "12345", "1234", "1234")
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
