package roomclient

import "time"

type RoomDTO struct {
	ID          uint     `json:"id"`
	HostID      uint     `json:"hostID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	MinGuests   uint     `json:"minGuests"`
	MaxGuests   uint     `json:"maxGuests"`
	Photos      []string `json:"photos"`
	Commodities []string `json:"commodities"`
	AutoApprove bool     `json:"autoApprove"`
	Deleted     bool     `json:"deleted"`
}

type CreateRoomDTO struct {
	HostID        uint     `json:"hostID"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Address       string   `json:"address"`
	MinGuests     uint     `json:"minGuests"`
	MaxGuests     uint     `json:"maxGuests"`
	PhotosPayload []string `json:"photosPayload"`
	Commodities   []string `json:"commodities"`
	AutoApprove   bool     `json:"autoApprove"`
}

// ---------------------------------------------------------------

type CreateRoomAvailabilityListDTO struct {
	RoomID uint                            `json:"roomId"`
	Items  []CreateRoomAvailabilityItemDTO `json:"items"`
}

type CreateRoomAvailabilityItemDTO struct {
	// ExistingID is either the ID of an RoomAvailabilityItem that already
	// exists, or 0 if this is a new item. When 0, a new one will be created in
	// the DB. When not 0, it will reuse the existing object.
	ExistingID uint      `json:"existingId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	Available  bool      `json:"available"`
}

// ---------------------------------------------------------------

type CreateRoomPriceListDTO struct {
	RoomID    uint                     `json:"roomId"`
	Items     []CreateRoomPriceItemDTO `json:"items"`
	BasePrice uint                     `json:"basePrice"`
	PerGuest  bool                     `json:"perGuest"`
}

type CreateRoomPriceItemDTO struct {
	// ExistingID is either the ID of an RoomPriceItem that already
	// exists, or 0 if this is a new item. When 0, a new one will be created in
	// the DB. When not 0, it will reuse the existing object.
	ExistingID uint      `json:"existingId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	Price      uint      `json:"price"`
}

// ---------------------------------------------------------------
