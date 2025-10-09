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
}

type ReservationDTO struct {
	ID                 uint      `gorm:"primaryKey"`
	RoomID             uint      `gorm:"not null"`
	RoomAvailabilityID uint      `gorm:"not null"`
	RoomPriceID        uint      `gorm:"not null"`
	GuestID            uint      `gorm:"not null"`
	DateFrom           time.Time `gorm:"not null"`
	DateTo             time.Time `gorm:"not null"`
	GuestCount         uint      `gorm:"not null"`
	Cancelled          bool      `gorm:"not null"`
	Cost               uint      `gorm:"not null"`
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
}
