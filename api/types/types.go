package types

import "time"

type ServiceStatistics struct {
	ServiceID string           `json:"service_id,omitempty"`
	CreatedAt time.Time        `json:"create_at,omitempty"`
	RemoveAt  time.Time        `json:"remove_at,omitempty"`
	State     string           `json:"state,omitempty"`
	TaskList  []TaskStatistics `json:"tasklist,omitempty"`
}
type TaskStatistics struct {
	TaskID         string    `json:"task_id,omitempty"`
	NodeID         string    `json:"node_id,omitempty"`
	ReceiveAddress string    `json:"receive_address,omitempty"`
	CreatedAt      time.Time `json:"create_at,omitempty"`
	RemoveAt       time.Time `json:"remove_at,omitempty"`
	State          string    `json:"state,omitempty"`
	Msg            string    `json:"msg,omitempty"`
	Err            string    `json:"err,omitempty"`
	DesiredState   string    `json:"desired_state,omitempty"`
}
