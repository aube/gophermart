package client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func ParseOrderAccrual(requestBody []byte) (OrderAccrual, error) {
	var oa OrderAccrual

	err := json.Unmarshal(requestBody, &oa)
	if err != nil {
		return OrderAccrual{}, err
	}

	return oa, nil
}

// User ...
type OrderAccrual struct {
	ID      int    `json:"order"`
	Accrual int    `json:"accrual"`
	Status  string `json:"status"`
}

// Request ...
func request(address string) (OrderAccrual, error) {
	resp, err := http.Get(address)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrderAccrual{}, err
	}

	switch resp.StatusCode {
	case 204:
		return OrderAccrual{}, errors.New("INVALID")
	case 429:
		return OrderAccrual{}, errors.New("NEW")
	default:
		return ParseOrderAccrual(body)
	}
}
