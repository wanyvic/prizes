package order

import "time"

type RefundPayment struct {
	RefundID          string    `json:"refund_id,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	TotalAmount       int64     `json:"total_amount,omitempty"`
	Payments          []Payment `json:"payments,omitempty"`
	RefundTransaction string    `json:"refund_transaction,omitempty"`
}
