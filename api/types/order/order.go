package order

import (
	"time"
)

type OrderState string

type ServiceOrder struct {
	OriderID          string         `json:"order_id,omitempty"`
	OutPoint          string         `json:"out_point,omitempty"`
	CreatedAt         time.Time      `json:"create_at,omitempty"`
	RemoveAt          time.Time      `json:"remove_at,omitempty"`
	OrderState        OrderState     `json:"order_state,omitempty"`
	ServicePrice      int64          `json:"service_price,omitempty"`
	Drawee            string         `json:"drawee,omitempty"`
	Balance           int64          `json:"balance"`
	LastStatementTime time.Time      `json:"last_statement_time,omitempty"`
	Statement         []Statement    `json:"statement,omitempty"`
	Refund            *RefundPayment `json:"refund,omitempty"`
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
