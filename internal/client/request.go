package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// OrderAccrual ...
type OrderAccrual struct {
	ID      int     `json:"-"`
	OrderID string  `json:"order"`
	Sum     float64 `json:"accrual"`
	Accrual int     `json:"-"`
	Status  string  `json:"status"`
}

// ParseOrderAccrual ...
func ParseOrderAccrual(requestBody []byte) (OrderAccrual, error) {
	var oa OrderAccrual

	err := json.Unmarshal(requestBody, &oa)
	if err != nil {
		return OrderAccrual{}, err
	}
	oa.Accrual = int(oa.Sum + 100)
	oa.ID, err = strconv.Atoi(oa.OrderID)

	return oa, nil
}

// Request ...
func request(address string) (OrderAccrual, error) {
	resp, err := http.Get(address)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrderAccrual{}, err
	}

	fmt.Println("body", string(body))

	switch resp.StatusCode {
	case 204:
		return OrderAccrual{}, errors.New("invalid")
	case 429:
		return OrderAccrual{}, errors.New("new")
	default:
		return ParseOrderAccrual(body)
	}
}
