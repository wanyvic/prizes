package service

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api/types/order"
)

//PrizesService includes order createSpec updateSpec and etc.
type PrizesService struct {
	DockerSerivce swarm.Service
	CreateSpec    ServiceCreate
	UpdateSpec    []ServiceUpdate
	Order         []order.ServiceOrder
	State         ServiceState
	CreatedAt     time.Time
	DeleteAt      time.Time
	NextCheckTime time.Time
}

//ServiceState includes UNKNOWN RUNNING COMPLETE
type ServiceState string

const (
	// ServiceStateUnknown UNKNOWN
	ServiceStateUnknown ServiceState = "unknown"
	// ServiceStateRunning RUNNING
	ServiceStateRunning ServiceState = "running"
	// ServiceStateCompleted COMPLETE
	ServiceStateCompleted ServiceState = "completed"
	//DefaultDockerImage massgrid/10.0-base-ubuntu16.04
	DefaultDockerImage = "massgrid/10.0-base-ubuntu16.04"
	// DefaultServiceCreateID 100000
	DefaultServiceCreateID = "100000"
	// DefaultServiceUpdateID 100100
	DefaultServiceUpdateID = "100100"
	// DefaultServiceRefundID 100200
	DefaultServiceRefundID = "100200"
	// DefaultStatementID 100300
	DefaultStatementID = "100300"
)

type ServiceInfo struct {
	ServiceID     string               `json:"service_id,omitempty"`
	CreatedAt     time.Time            `json:"create_at,omitempty"`
	DeleteAt      time.Time            `json:"delete_at,omitempty"`
	NextCheckTime time.Time            `json:"next_check_time,omitempty"`
	Order         []order.ServiceOrder `json:"order,omitempty"`
	CreateSpec    ServiceCreate        `json:"create_spec,omitempty"`
	UpdateSpec    []ServiceUpdate      `json:"update_spec,omitempty"`
	State         ServiceState         `json:"state,omitempty"`
	TaskInfo      *swarm.Task          `json:"task_info,omitempty"`
}
