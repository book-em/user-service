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
	DeleteHostRooms(ctx context.Context, jwt string) ([]RoomDTO, error)
}

type roomClient struct {
	baseURL string
}

func NewRoomClient() RoomClient {
	return &roomClient{
		baseURL: "http://room-service:8080/api/", // TODO: This should not be hardcoded
	}
}

func (c *roomClient) DeleteHostRooms(ctx context.Context, jwt string) ([]RoomDTO, error) {
	utils.TEL.Push(ctx, "delete-host-rooms")
	defer utils.TEL.Pop()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%shost/", c.baseURL), nil)
	if err != nil {
		utils.TEL.Error("preparing request error ", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		utils.TEL.Error("request error ", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TEL.Error("parsing response error", err)
		return nil, err
	}

	var obj []RoomDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		utils.TEL.Error("JSON unmarshall error", err)
		return nil, err
	}

	return obj, nil
}
