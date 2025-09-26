package reservationclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ReservationClient interface {
	// GetActiveGuestReservations finds all reservations made by `guest` that
	// haven't completed yet. The user must be a guest.
	GetActiveGuestReservations(jwt string) ([]ReservationDTO, error)
}

type reservationClient struct {
	baseURL string
}

func NewReservationClient() ReservationClient {
	return &reservationClient{
		baseURL: "http://reservation-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *reservationClient) GetActiveGuestReservations(jwt string) ([]ReservationDTO, error) {
	log.Printf("Get active guest reservations")

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/reservations/guest/active", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error %v", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Parsing response error: %v", err)
		return nil, err
	}

	var obj []ReservationDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		log.Printf("JSON Unmarshall error: %v", err)
		return nil, err
	}

	return obj, nil
}
