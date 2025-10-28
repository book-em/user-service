package reservationclient

import "time"

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

type CreateReservationRequestDTO struct {
	RoomID     uint      `json:"roomId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	GuestCount uint      `json:"guestCount"`
}

type ReservationRequestDTO struct {
	ID               uint      `json:"id"`
	RoomID           uint      `json:"roomId"`
	DateFrom         time.Time `json:"dateFrom"`
	DateTo           time.Time `json:"dateTo"`
	GuestCount       uint      `json:"guestCount"`
	GuestID          uint      `json:"guestId"`
	Status           string    `json:"status"`
	Cost             uint      `json:"cost"`
	GuestCancelCount uint      `json:"guestCancelCount"`
}
