package order

import "time"

type RefundInfo struct {
	RefundPay         *[]RefundPayment `json:refund_pay,omitempty""`
	Statement         *Statement       `json:statement,omitempty""`
	RefundTransaction string           `json:"refund_transaction,omitempty"`
}
type RefundPayment struct {
	RefundID    string    `json:"refund_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	TotalAmount int64     `json:"total_amount,omitempty"`
	Drawee      string    `json:"drawee,omitempty"`
}
