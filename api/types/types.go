package types

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type ServiceStatistics struct {
	ServiceID string          `json:"service_id,omitempty"`
	CreatedAt time.Time       `json:"create_at,omitempty"`
	RemoveAt  time.Time       `json:"remove_at,omitempty"`
	State     swarm.TaskState `json:"state,omitempty"`
	Price     int64           `json:"price,omitempty"`
	FeeRate   float64         `json:"fee_rate,omitempty"`
	Hardware
	TaskList []TaskStatistics `json:"tasklist,omitempty"`
}
type TaskStatistics struct {
	TaskID         string          `json:"task_id,omitempty"`
	NodeID         string          `json:"node_id,omitempty"`
	ReceiveAddress string          `json:"receive_address,omitempty"`
	CreatedAt      time.Time       `json:"create_at,omitempty"`
	RemoveAt       time.Time       `json:"remove_at,omitempty"`
	State          swarm.TaskState `json:"state,omitempty"`
	Msg            string          `json:"msg,omitempty"`
	Err            string          `json:"err,omitempty"`
	DesiredState   swarm.TaskState `json:"desired_state,omitempty"`
}
