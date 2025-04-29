package client

import (
	"fmt"
	"strconv"
)

type OrderProvider interface {
	SetStatus(int, string) error
	SetAccrual(int, int) error
}

type accrualClient struct {
	systemAddress string
	storeOrder    OrderProvider
}

func New(systemAddress string, storeOrder OrderProvider) *accrualClient {
	return &accrualClient{
		systemAddress: systemAddress,
		storeOrder:    storeOrder,
	}
}

func (ac *accrualClient) SendOrder(id int) string {
	fmt.Println("Accrual order", id)

	oa, err := request(ac.systemAddress + "/api/orders/" + strconv.Itoa(id))

	fmt.Println("Accrual service response", oa, err)

	if err.Error() == "new" {
		ac.storeOrder.SetStatus(id, "NEW")
		return ""
	}

	if err.Error() == "invalid" {
		ac.storeOrder.SetStatus(id, "INVALID")
		return ""
	}

	if oa.Status == "INVALID" {
		ac.storeOrder.SetStatus(id, "INVALID")
		return ""
	}

	ac.storeOrder.SetAccrual(id, oa.Accrual)
	return ""
}
