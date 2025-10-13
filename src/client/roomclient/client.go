package roomclient

type RoomClient interface {
}

type roomClient struct {
	baseURL string
}

func NewRoomClient() RoomClient {
	return &roomClient{
		baseURL: "http://room-service:8080/api", // TODO: This should not be hardcoded
	}
}
