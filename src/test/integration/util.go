package test

import (
	"bookem-user-service/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
