package types

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type ServiceStatistics struct {
	ServiceID string           `json:"service_id,omitempty"`
	StartAt   time.Time        `json:"start_at,omitempty"`
	EndAt     time.Time        `json:"end_at,omitempty"`
	TaskList  []TaskStatistics `json:"tasklist,omitempty"`
}
type TaskStatistics struct {
	TaskID         string          `json:"task_id,omitempty"`
	NodeID         string          `json:"node_id,omitempty"`
	ReceiveAddress string          `json:"receive_address,omitempty"`
	StartAt        time.Time       `json:"start_at,omitempty"`
	EndAt          time.Time       `json:"end_at,omitempty"`
	State          swarm.TaskState `json:"state,omitempty"`
	Msg            string          `json:"msg,omitempty"`
	Err            string          `json:"err,omitempty"`
	DesiredState   swarm.TaskState `json:"desired_state,omitempty"`
}
