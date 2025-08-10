package roomclient

import "bookem-user-service/domain"

type RoomClient interface {
	// GetPendingGuestReservations finds all reservations made by `guest` that
	// haven't completed yet. The user must be a guest.
	GetPendingGuestReservations(guest *domain.User) ([]ReservationDTO, error)
	// GetActiveHostReservations finds all reservations made to rooms owned by
	// `host` that haven't completed yet. The user must be a host.
	GetActiveHostReservations(host *domain.User) ([]ReservationDTO, error)
}

type roomClient struct {
	baseURL string
}

func NewRoomClient() RoomClient {
	return &roomClient{
		baseURL: "http://localhost:9999", // Placeholder URL for now
	}
}

func (c *roomClient) GetPendingGuestReservations(guest *domain.User) ([]ReservationDTO, error) {
	if guest.Role != domain.Guest {
		return []ReservationDTO{}, domain.ErrUnauthorized
	}

	return []ReservationDTO{}, nil
}

func (c *roomClient) GetActiveHostReservations(host *domain.User) ([]ReservationDTO, error) {
	if host.Role != domain.Host {
		return []ReservationDTO{}, domain.ErrUnauthorized
	}

	return []ReservationDTO{}, nil
}
