package types

import "time"

type ServiceStatistics struct {
	ServiceID string
	CreatedAt time.Time
	RemoveAt  time.Time
	State     string
	TaskList  []TaskStatistics
}
type TaskStatistics struct {
	TaskID         string
	NodeID         string
	ReceiveAddress string
	CreatedAt      time.Time
	RemoveAt       time.Time
	State          string
	Msg            string
	Err            string
	DesiredState   string
}
