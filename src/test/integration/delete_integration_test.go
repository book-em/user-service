package test

import (
	"bookem-user-service/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Delete(t *testing.T) {
	{
		resp, _ := registerUser("guest1", "1234", domain.Guest)
		jwt := loginUser2("guest1", "1234")
		resp, err := deleteUser(jwt)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
	{
		resp, _ := registerUser("host1", "1234", domain.Host)
		jwt := loginUser2("host1", "1234")
		resp, err := deleteUser(jwt)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}

func TestIntegration_Delete_GuestHasActiveReservations(t *testing.T) {
	resp, _ := registerUser("guest2", "1234", domain.Guest)
	jwt := loginUser2("guest2", "1234")

	// TODO: Once we can actually create reservations
	//       we  will  add  active reservations here.
	//       Until  then,  this   test   will   pass.

	resp, err := deleteUser(jwt)
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_GuestHasNoActiveReservations(t *testing.T) {
	resp, _ := registerUser("guest3", "1234", domain.Guest)
	jwt := loginUser2("guest3", "1234")

	// TODO: Once we can actually create reservations
	//       we will add non-active reservations here.
	//       Until  then,  this   test   will   pass.

	resp, err := deleteUser(jwt)
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_HostHasActiveReservations(t *testing.T) {
	resp, _ := registerUser("host2", "1234", domain.Host)
	id := getUserFromRegister(resp).Id
	jwt := loginUser2("host2", "1234")

	roomDTO := DefaultRoomCreateDTO
	roomDTO.HostID = id

	resp, _ = createRoom(jwt, roomDTO)

	// TODO: Once we can actually create reservations
	//       we  will  add  active reservations here.
	//       Until  then,  this   test   will   pass.

	resp, err := deleteUser(jwt)
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_UserIsAdmin(t *testing.T) {
	resp, _ := registerUser("admin1", "1234", domain.Admin)
	jwt := loginUser2("admin1", "1234")
	resp, err := deleteUser(jwt)

	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
