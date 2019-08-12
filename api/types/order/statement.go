package order

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type Statement struct {
	StatementID          string    `json:"statement_id,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	StatementStartAt     time.Time `json:"statement_start_at,omitempty"`
	StatementEndAt       time.Time `json:"statement_end_at,omitempty"`
	TotalAmount          int64     `json:"total_amount,omitempty"`
	MasterNodeFeeRate    int64     `json:"master_node_fee_rate,omitempty"`
	DevFeeRate           int64     `json:"dev_fee_rate,omitempty"` //max 10000
	MasterNodeFeeAddress string    `json:"master_node_fee_address,omitempty"`
	DevFeeAddress        string    `json:"dev_fee_address,omitempty"`
	StatementTransaction string    `json:"statement_transaction,omitempty"`
	Payments             []Payment `json:"payments,omitempty"`
}
type Payment struct {
	ReceiveAddress string          `json:"receive_address,omitempty"`
	Amount         int64           `json:"amount,omitempty"`
	TaskID         string          `json:"task_id,omitempty"`
	TaskState      swarm.TaskState `json:"task_state,omitempty"`
	Msg            string          `json:"msg,omitempty"`
}

type StatementOptions struct {
	MasterNodeFeeRate    int64  `json:"master_node_fee_rate,omitempty"`
	DevFeeRate           int64  `json:"dev_fee_rate,omitempty"`
	MasterNodeFeeAddress string `json:"master_node_fee_address,omitempty"`
	DevFeeAddress        string `json:"dev_fee_address,omitempty"`
}
