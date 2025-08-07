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
