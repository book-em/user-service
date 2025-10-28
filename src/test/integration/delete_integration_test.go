package test

import (
	"bookem-user-service/client/reservationclient"
	"bookem-user-service/domain"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Delete(t *testing.T) {
	{
		resp, _ := registerUser("guest_idel_1", "1234", domain.Guest)
		jwt := loginUser2("guest_idel_1", "1234")
		resp, err := deleteUser(jwt)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
	{
		resp, _ := registerUser("host_idel_1", "1234", domain.Host)
		jwt := loginUser2("host_idel_1", "1234")
		resp, err := deleteUser(jwt)

		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}

func TestIntegration_Delete_GuestHasActiveReservations(t *testing.T) {
	hostUsername := "host_idel_2"
	_, _, hostJwt, room := setupHostRoomAvailabilityPrice(hostUsername, t)

	registerUser("guest_idel_2", "1234", domain.Guest)
	guestJwt := loginUser2("guest_idel_2", "1234")

	dto := reservationclient.CreateReservationRequestDTO{
		RoomID:     room.ID,
		DateFrom:   time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	resp, err := createReservationRequest(guestJwt, dto)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	req := responseToReservationRequest(resp)

	approveURL := URL_reservation + "req/" + strconv.FormatUint(uint64(req.ID), 10) + "/approve"
	request, err := http.NewRequest(http.MethodPut, approveURL, nil)
	require.NoError(t, err)
	request.Header.Add("Authorization", "Bearer "+hostJwt)

	approveResp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, approveResp.StatusCode)

	resp, err = deleteUser(hostJwt)
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_GuestHasNoActiveReservations(t *testing.T) {
	resp, _ := registerUser("guest_idel_3", "1234", domain.Guest)
	guestJwt := loginUser2("guest_idel_3", "1234")

	hostUsername := "host_idel_3"
	_, _, _, room := setupHostRoomAvailabilityPrice(hostUsername, t)

	// This reservation request is inactive.
	dto := reservationclient.CreateReservationRequestDTO{
		RoomID:     room.ID,
		DateFrom:   time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	resp, err := createReservationRequest(guestJwt, dto)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = deleteUser(guestJwt)
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegration_Delete_HostHasActiveReservations(t *testing.T) {
	hostUsername := "host_idel_4"
	_, _, hostJwt, room := setupHostRoomAvailabilityPrice(hostUsername, t)

	registerUser("guest_idel_4", "1234", domain.Guest)
	guestJwt := loginUser2("guest_idel_4", "1234")

	dto := reservationclient.CreateReservationRequestDTO{
		RoomID:     room.ID,
		DateFrom:   time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	resp, err := createReservationRequest(guestJwt, dto)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	req := responseToReservationRequest(resp)

	approveURL := URL_reservation + "req/" + strconv.FormatUint(uint64(req.ID), 10) + "/approve"
	request, err := http.NewRequest(http.MethodPut, approveURL, nil)
	require.NoError(t, err)
	request.Header.Add("Authorization", "Bearer "+hostJwt)

	approveResp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, approveResp.StatusCode)

	resp, err = deleteUser(hostJwt)
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
