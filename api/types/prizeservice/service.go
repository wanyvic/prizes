package prizeservice

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type PrizesService struct {
	DockerSerivce swarm.Service
	CreateSpec    ServiceCreate
	Order         []ServiceOrder
	Record        []Statement
	State         string
	CreatedAt     time.Time
	DeleteAt      time.Time
}

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
