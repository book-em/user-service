package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Delete(t *testing.T) {
	{
		resp, _ := registerUser("user_deleted_guest_01", "1234", domain.Guest)
		id := getUserFromRegister(resp).Id
		jwt := loginUser2("user_deleted_guest_01", "1234")
		resp, err := deleteUserById(jwt, id)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
	{
		resp, _ := registerUser("user_deleted_host_01", "1234", domain.Host)
		id := getUserFromRegister(resp).Id
		jwt := loginUser2("user_deleted_host_01", "1234")
		resp, err := deleteUserById(jwt, id)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}

func TestIntegration_Delete_WrongUser(t *testing.T) {
	registerUser("user_deleted_guest_03", "1234", domain.Guest)
	id := uint(999999)
	jwt := loginUser2("user_deleted_guest_03", "1234")
	resp, err := deleteUserById(jwt, id)

	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_Delete_GuestHasPendingReservations(t *testing.T) {
	resp, _ := registerUser("user_deleted_guest_02", "1234", domain.Guest)
	id := getUserFromRegister(resp).Id
	jwt := loginUser2("user_deleted_guest_02", "1234")
	resp, err := deleteUserById(jwt, id)

	// TODO: Once we can actually create reservations, that needs to happen here so we can trigger a 400.
	// Until then, this test will pass.

	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_HostHasPendingReservations(t *testing.T) {
	resp, _ := registerUser("user_deleted_host_02", "1234", domain.Host)
	id := getUserFromRegister(resp).Id
	jwt := loginUser2("user_deleted_host_02", "1234")
	resp, err := deleteUserById(jwt, id)

	// TODO: Once we can actually create reservations, that needs to happen here so we can trigger a 400.
	// Until then, this test will pass.

	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_UserIsAdmin(t *testing.T) {
	resp, _ := registerUser("user_deleted_admin_01", "1234", domain.Admin)
	id := getUserFromRegister(resp).Id
	jwt := loginUser2("user_deleted_admin_01", "1234")
	resp, err := deleteUserById(jwt, id)

	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
