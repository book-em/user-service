package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Update_Success(t *testing.T) {
	resp, err := registerUser("user_10", "1234", domain.Guest)
	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_10", "1234")

	new_username := "user_10_NEW"

	resp, err = updateUser(jwt, user.Id, &new_username, nil)
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Update_WrongUser(t *testing.T) {
	registerUser("user_11_somebody_else", "1234", domain.Guest)
	resp, _ := registerUser("user_11", "1234", domain.Guest)

	user := getUserFromRegister(resp)
	new_username := "user_11_NEW"

	// Somebody else.

	jwt := loginUser2("user_11_somebody_else", "1234")
	resp, err := updateUser(jwt, user.Id, &new_username, nil)
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Signed out user.

	resp, err = updateUser("", user.Id, &new_username, nil)
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestIntegration_Update_UsernameTaken(t *testing.T) {
	resp, _ := registerUser("user_13_1", "1234", domain.Guest)
	registerUser("user_13_2", "1234", domain.Guest)

	// Change user_13_1's username to user_13_2's username

	user := getUserFromRegister(resp)
	jwt := loginUser2("user_13_1", "1234")
	new_username := "user_13_2"

	resp, err := updateUser(jwt, user.Id, &new_username, nil)
	if err != nil {
		t.Log(err.Error())
	}
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
}
