package roomclient

import (
	"bookem-user-service/domain"
	utils "bookem-user-service/util"
	"context"
)

type RoomClient interface {
	// GetPendingGuestReservations finds all reservations made by `guest` that
	// haven't completed yet. The user must be a guest.
	GetPendingGuestReservations(ctx context.Context, guest *domain.User) ([]ReservationDTO, error)
	// GetActiveHostReservations finds all reservations made to rooms owned by
	// `host` that haven't completed yet. The user must be a host.
	GetActiveHostReservations(ctx context.Context, host *domain.User) ([]ReservationDTO, error)
}

type roomClient struct {
	baseURL string
}

func NewRoomClient() RoomClient {
	return &roomClient{
		baseURL: "http://localhost:9999", // Placeholder URL for now
	}
}

func (c *roomClient) GetPendingGuestReservations(ctx context.Context, guest *domain.User) ([]ReservationDTO, error) {
	utils.TEL.Push(ctx, "get-reservation-requests-made-by-guest")
	defer utils.TEL.Pop()

	if guest.Role != domain.Guest {
		return []ReservationDTO{}, domain.ErrUnauthorized
	}

	return []ReservationDTO{}, nil
}

func (c *roomClient) GetActiveHostReservations(ctx context.Context, host *domain.User) ([]ReservationDTO, error) {
	utils.TEL.Push(ctx, "get-reservation-requests-for-host")
	defer utils.TEL.Pop()

	if host.Role != domain.Host {
		return []ReservationDTO{}, domain.ErrUnauthorized
	}

	return []ReservationDTO{}, nil
}
