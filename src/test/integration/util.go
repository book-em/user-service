package test

import (
	"bookem-user-service/domain"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
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
	body := resp.Body.Close().Error()
	var token domain.JWTDTO
	json.Unmarshal([]byte(body), &token)
	return token.Jwt
}
