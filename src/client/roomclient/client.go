package roomclient

import (
	utils "bookem-user-service/util"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RoomClient interface {
	// GetActiveHostReservations finds all reservations made to rooms owned by
	// `host` that haven't completed yet. The user must be a host.
	GetActiveHostReservations(ctx context.Context, jwt string) ([]ReservationDTO, error)
}

type roomClient struct {
	baseURL string
}

func NewRoomClient() RoomClient {
	return &roomClient{
		baseURL: "http://room-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *roomClient) GetActiveHostReservations(ctx context.Context, jwt string) ([]ReservationDTO, error) {
	utils.TEL.Push(ctx, "get-active-reservations-for-host")
	defer utils.TEL.Pop()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/reservations/host/active", c.baseURL), nil)
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
