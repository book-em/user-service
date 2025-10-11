package reservationclient

import (
	utils "bookem-user-service/util"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReservationClient interface {
	// GetActiveGuestReservations finds all reservations made by `guest` that
	// haven't completed yet. The user must be a guest.
	GetActiveGuestReservations(ctx context.Context, jwt string) ([]ReservationDTO, error)
}

type reservationClient struct {
	baseURL string
}

func NewReservationClient() ReservationClient {
	return &reservationClient{
		baseURL: "http://reservation-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *reservationClient) GetActiveGuestReservations(ctx context.Context, jwt string) ([]ReservationDTO, error) {
	utils.TEL.Push(ctx, "get-active-reservations-for-guest")
	defer utils.TEL.Pop()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/reservations/guest/active", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		utils.TEL.Error("error ", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TEL.Error("parsing response error", err)
		return nil, err
	}

	var obj []ReservationDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		utils.TEL.Error("JSON unmarshall error", err)
		return nil, err
	}

	return obj, nil
}
