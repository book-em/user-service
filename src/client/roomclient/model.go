package roomclient

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
