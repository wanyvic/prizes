package service

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/wanyvic/prizes/api/types/order"
)

type PrizesService struct {
	DockerSerivce swarm.Service
	CreateSpec    ServiceCreate
	UpdateSpec    []ServiceUpdate
	Order         []order.ServiceOrder
	Record        []order.Statement
	State         ServiceState
	CreatedAt     time.Time
	DeleteAt      time.Time
	NextCheckTime time.Time
}
type ServiceState string

const (
	// NodeStateUnknown UNKNOWN
	ServiceStateUnknown ServiceState = "unknown"
	// NodeStateReady READY
	ServiceStateRunning ServiceState = "running"
	// NodeStateDisconnected DISCONNECTED
	ServiceStateCompleted ServiceState = "completed"
)

var (
	DefaultServiceCreateID = "100000"
	DefaultServiceUpdateID = "100100"
	DefaultServiceRefundID = "100200"
	DefaultStatementID     = "100300"
)

// 创建 服务
// 通过 serviceCreate 配置信息创建服务 返回 serviceID 和错误信息
type ServiceCommand interface {
	// ServiceCreate(*prizestypes.ServiceCreate) (string, error)

	// ServiceReFund() (prizestypes.Statement, error)
}

//退款服务
// 返回结算清单
//结算服务

//更新服务

//
