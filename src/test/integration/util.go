package test

import (
	middleware "bookem-user-service/api/middleware"
	"bookem-user-service/client/reservationclient"
	"bookem-user-service/client/roomclient"
	"bookem-user-service/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const URL = "http://user-service:8080/api/"
const URL_room = "http://room-service:8080/api/"
const URL_reservation = "http://reservation-service:8080/api/"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func registerUser(username_or_email string, password string, role domain.UserRole) (*http.Response, error) {
	username := username_or_email
	email := username + "@gmail.com"

	if strings.HasSuffix(username_or_email, "@gmail.com") {
		username = strings.Split(username_or_email, "@")[0]
		email = username_or_email
	}

	dto := domain.UserCreateDTO{
		Username: username,
		Password: password,
		Email:    email,
		Role:     string(role),
		Name:     genName(6),
		Surname:  genName(6),
		Address:  genName(10),
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL+"register", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func getUserFromRegister(resp *http.Response) domain.UserDTO {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var user domain.UserDTO
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		panic(fmt.Sprintf("failed to unmarshal user: %v", err))
	}

	return user
}

func loginUser(username_or_email string, password string) (*http.Response, error) {
	dto := domain.LoginDTO{
		UsernameOrEmail: username_or_email,
		Password:        password,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL+"login", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func loginUser2(username_or_email string, password string) string {

	resp, _ := loginUser(username_or_email, password)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var token domain.JWTDTO
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		panic(fmt.Sprintf("failed to unmarshal jwt: %v", err))
	}

	return token.Jwt
}

func updateUser(jwt string, id uint, new_username *string, new_surname *string) (*http.Response, error) {
	dto := domain.UserUpdateDTO{
		Id:       id,
		Username: new_username,
		Surname:  new_surname,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, URL+"update", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func changePassword(jwt string, id uint, old, new, newConfirm string) (*http.Response, error) {
	dto := domain.PasswordUpdateDTO{
		Id:                 id,
		OldPassword:        old,
		NewPassword:        new,
		NewPasswordConfirm: newConfirm,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, URL+"password", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func findUserById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s%d", URL, id)) // No forward slash between them, it's in `URL`
	return resp, err
}

func deleteUser(jwt string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func createRoom(jwt string, dto roomclient.CreateRoomDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"new", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func createReservationRequest(jwt string, dto reservationclient.CreateReservationRequestDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_reservation+"req", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func responseToReservationRequest(resp *http.Response) reservationclient.ReservationRequestDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj reservationclient.ReservationRequestDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func responseToRoom(resp *http.Response) roomclient.RoomDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj roomclient.RoomDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		fmt.Print(string(bodyBytes))
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func createRoomAvailability(jwt string, dto roomclient.CreateRoomAvailabilityListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"available", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func createRoomPrice(jwt string, dto roomclient.CreateRoomPriceListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"price", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func setupHostRoomAvailabilityPrice(hostUsername string, t *testing.T) (string, string, string, roomclient.RoomDTO) {
	// Step 1: Register unique host
	username := hostUsername
	password := "pass"
	registerUser(username, password, domain.Host)
	jwt := loginUser2(username, password)
	jwtObj, err := middleware.GetJwtFromString(jwt)
	require.NoError(t, err)

	// Step 2: Create room
	roomDTO := roomclient.CreateRoomDTO{
		HostID:        jwtObj.ID,
		Name:          "Room_" + genName(6),
		Description:   "Test room",
		Address:       "Test address",
		MinGuests:     1,
		MaxGuests:     4,
		PhotosPayload: []string{SMALL_IMG},
		Commodities:   []string{"WiFi", "AC"},
		AutoApprove:   false,
	}
	roomResp, err := createRoom(jwt, roomDTO)
	require.NoError(t, err)
	defer roomResp.Body.Close()
	room := responseToRoom(roomResp)

	// Step 3: Create availability list
	availabilityDTO := roomclient.CreateRoomAvailabilityListDTO{
		RoomID: room.ID,
		Items: []roomclient.CreateRoomAvailabilityItemDTO{
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
				Available:  true,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC),
				Available:  true,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
				Available:  true,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
				Available:  true,
			},
		},
	}
	availResp, err := createRoomAvailability(jwt, availabilityDTO)
	require.NoError(t, err)
	defer availResp.Body.Close()

	// Step 4: Create price list
	priceDTO := roomclient.CreateRoomPriceListDTO{
		RoomID:    room.ID,
		BasePrice: 80,
		PerGuest:  false,
		Items: []roomclient.CreateRoomPriceItemDTO{
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
				Price:      100,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC),
				Price:      120,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
				Price:      200,
			},
			{
				ExistingID: 0,
				DateFrom:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				DateTo:     time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
				Price:      200,
			},
		},
	}
	priceResp, err := createRoomPrice(jwt, priceDTO)
	require.NoError(t, err)
	defer priceResp.Body.Close()

	return username, password, jwt, room
}

// ----------------------------------------------- Mock data

const (
	SMALL_IMG = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQQEBQoHBwYIDAoMDAsKCwsNDhIQDQ4RDgsLEBYQERMUFRUVDA8XGBYUGBIUFRT/wAALCAABAAEBAREA/8QAFAABAAAAAAAAAAAAAAAAAAAACf/EABQQAQAAAAAAAAAAAAAAAAAAAAD/2gAIAQEAAD8AKp//2Q=="
)

var DefaultRoomCreateDTO = roomclient.CreateRoomDTO{
	HostID:        1,
	Name:          "Room Name",
	Description:   "Room Desc",
	Address:       "Room Address",
	MinGuests:     1,
	MaxGuests:     5,
	PhotosPayload: []string{SMALL_IMG},
	Commodities:   []string{"WiFi"},
}
