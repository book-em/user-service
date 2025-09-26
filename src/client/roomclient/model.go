package roomclient

import "time"

type RoomDTO struct {
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
