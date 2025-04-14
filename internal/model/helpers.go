package model

import (
	"encoding/json"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ParseCredentials(requestBody []byte) (User, error) {
	var user User

	err := json.Unmarshal(requestBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ParseOrderID(requestBody []byte) (Order, error) {
	id, err := strconv.Atoi(string(requestBody))
	if err != nil {
		return Order{}, err
	}

	return Order{
		ID:     id,
		UserID: "33",
	}, nil
}

func OrdersToJSON(orders []Order) ([]byte, error) {
	jsonData, err := json.Marshal(orders)

	return jsonData, err
}
