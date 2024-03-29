package order

import (
	"time"
)

type OrderState string

type ServiceOrder struct {
	OriderID              string        `json:"order_id,omitempty"`
	OutPoint              string        `json:"out_point,omitempty"`
	CreatedAt             time.Time     `json:"create_at,omitempty"`
	OrderState            OrderState    `json:"order_state,omitempty"`
	ServicePrice          int64         `json:"service_price,omitempty"`
	MasterNodeFeeRate     int64         `json:"master_node_fee_rate"`
	DevFeeRate            int64         `json:"dev_fee_rate"` //max 10000
	MasterNodeFeeAddress  string        `json:"master_node_fee_address,omitempty"`
	DevFeeAddress         string        `json:"dev_fee_address,omitempty"`
	Drawee                string        `json:"drawee,omitempty"`
	Balance               int64         `json:"balance"`
	Refund                int64         `json:"refund"`
	PayAmount             int64         `json:"pay_amount"`
	RemainingTimeDuration time.Duration `json:"remaining_time_duration"`
	TotalTimeDuration     time.Duration `json:"total_time_duration"`
	LastStatementTime     time.Time     `json:"last_statement_time,omitempty"`
	Statement             []Statement   `json:"statement,omitempty"`
}

const (
	// OrderStateUnknown UNKNOWN
	OrderStateUnknown OrderState = "unknown"
	// OrderStateWaitToPay WAITTING
	OrderStateWaitToPay OrderState = "waitting"
	// OrderStatePaying PAYING
	OrderStatePaying OrderState = "paying"
	// OrderStateHasBeenPaid PAID
	OrderStateHasBeenPaid OrderState = "paid"
	// OrderStateHasBeenRefund REFUND
	OrderStateHasBeenRefund OrderState = "refund"
)
