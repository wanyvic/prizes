package prizeservice

import "time"

type ServiceOrder struct {
	OutPoint   string    `json:"order_id,omitempty"`
	CreatedAt  time.Time `json:"create_at,omitempty"`
	RemoveAt   time.Time `json:"create_at,omitempty"`
	OrderState string    `json:"order_state,omitempty"`

	Statement         []Statement
	LastStatementTime time.Time
	Balance           int64
}

type Statement struct {
	StatementID          string    `json:"statement_id,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	TotalAmount          int64     `json:"total_amount,omitempty"`
	MasterNodeFeeRate    float64   `json:"master_node_fee_rate,omitempty"`
	DevFeeRate           float64   `json:"dev_fee_rate,omitempty"`
	MasterNodeFeeAddress string    `json:"master_node_fee_address,omitempty"`
	DevFeeAddress        string    `json:"dev_fee_address,omitempty"`
	Payments             []Payment `json:"payments,omitempty"`
}
type Payment struct {
	ReceiveAddress string `json:"receive_address,omitempty"`
	Amount         int64  `json:"amount,omitempty"`
	TaskID         string `json:"task_id,omitempty"`
	TaskState      string `json:"task_state,omitempty"`
}
