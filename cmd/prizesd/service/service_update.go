package prizeservice

import (
	"time"

	"github.com/docker/docker/api/types"
	prizestypes "github.com/wanyvic/prizes/api/types"
)

type ServiceUpdate struct {
	ServiceID    string `json:"service_id,omitempty"`
	Amount       int64  `json:"amount,omitempty"`
	Pubkey       string `json:"pubkey,omitempty"`
	BlockHeight  int64  `json:"block_height,omitempty"`
	ServicePrice int64  `json:"service_price,omitempty"`
	OutPoint     string `json:"out_point,omitempty"`
}

func preaseUpdateServiceOrder(serviceUpdate *prizestypes.ServiceUpdate, response *types.ServiceUpdateResponse) (serviceOrder *prizestypes.ServiceOrder) {

	serviceOrder.OrderID = serviceUpdate.OutPoint
	serviceOrder.ServiceID = serviceUpdate.ServiceID
	serviceOrder.Msg = response.Warnings
	serviceOrder.BlockHeight = serviceUpdate.BlockHeight
	serviceOrder.ServicePrice = serviceUpdate.ServicePrice
	serviceOrder.CreatedAt = time.Now().UTC()
	serviceOrder.OrderState = "running"
	serviceOrder.OrderState = "running"
	return serviceOrder
}
