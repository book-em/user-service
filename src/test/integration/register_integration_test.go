package test

import (
	"bookem-user-service/domain"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ROLE_GUEST = "guest"
	ROLE_HOST  = "host"
	ROLE_ADMIN = "admin"
)

const URL = "http://user-service:8080/api/"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func registerUser(username_or_email string, password string, role string) (*http.Response, error) {
	username := username_or_email
	email := username + "@gmail.com"

	if strings.HasSuffix(username_or_email, "@gmail.com") {
		username = strings.Split(username_or_email, "@")[0]
		email = username_or_email
	}

	dto := domain.UserDTO{
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
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

func TestIntegration_RegisterUser(t *testing.T) {
	resp, err := registerUser("user1", "1234", ROLE_GUEST)
	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

// func TestIntegration_RegisterUserDuplicate(t *testing.T) {
// 	resp, err := registerUser("user1", "1234", ROLE_GUEST)
// 	require.Nil(t, err)
// 	require.Equal(t, resp.StatusCode, http.StatusCreated)

// 	resp, err = registerUser("user1", "1234", ROLE_GUEST)
// 	require.NotNil(t, err)
// 	require.Equal(t, resp.StatusCode, http.StatusConflict)

// 	resp, err = registerUser("user1@gmail.com", "1234", ROLE_GUEST)
// 	require.NotNil(t, err)
// 	require.Equal(t, resp.StatusCode, http.StatusConflict)
// }
