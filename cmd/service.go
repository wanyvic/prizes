package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/prizeservice"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
)

func ServiceCreate(serviceCreate *service.ServiceCreate) (*types.ServiceCreateResponse, error) {
	prizeService, response, err := prizeservice.Create(serviceCreate)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.Push(prizeService.DockerSerivce.ID, prizeService.NextCheckTime)
	return response, nil
}
func ServiceUpdate(serviceUpdate *service.ServiceUpdate) (*types.ServiceUpdateResponse, error) {
	prizeService, err := db.DBimplement.FindPrizesServiceOne(serviceUpdate.ServiceID)
	if err != nil {
		return nil, err
	}
	response, err := prizeservice.Update(prizeService, serviceUpdate)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.ChangeCheckTime(prizeService.DockerSerivce.ID, prizeService.NextCheckTime)
	return response, nil
}
func ServiceStatement(ServiceID string, statementAt time.Time) (*order.Statement, error) {
	prizeService, err := db.DBimplement.FindPrizesServiceOne(ServiceID)
	if err != nil {
		return nil, err
	}
	serviceStatistics, err := ServiceState(ServiceID, statementAt)
	if err != nil {
		return nil, err
	}
	statement, serviceState, err := prizeservice.Statement(prizeService, serviceStatistics, statementAt, order.DefaultStatementOptions)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService) //获取下次 结算时间
	if err != nil {
		return nil, err
	}
	if serviceState == service.ServiceStateRunning {
		calculagraph.ChangeCheckTime(prizeService.DockerSerivce.ID, prizeService.NextCheckTime)
	} else if serviceState == service.ServiceStateCompleted {
		err := serviceRemove(ServiceID)
		calculagraph.RemoveService(ServiceID)
		if err != nil {
			return nil, err
		}
	}
	logrus.Info(fmt.Sprintf("%+v", *statement))
	return statement, nil
}
func ServiceRefund(ServiceID string) (*order.RefundInfo, error) {
	var err error
	refundInfo := order.RefundInfo{}
	refundInfo.Statement, err = ServiceStatement(ServiceID, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	prizeService, err := db.DBimplement.FindPrizesServiceOne(ServiceID)
	if err != nil {
		return nil, err
	}
	refundInfo.RefundPay = prizeservice.Refund(prizeService)
	err = serviceRemove(ServiceID)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.RemoveService(ServiceID)

	logrus.Info(fmt.Sprintf("%+v", refundInfo))
	return &refundInfo, nil
}
func serviceRemove(serviceID string) error {
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	err = cli.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return err
	}
	err = refresh.RefreshStopService(serviceID)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("RemoveService completed: ID: %s", serviceID))
	return nil
}

func Service(serviceID string) (*swarm.Service, error) {
	logrus.Info("request ServiceInfo: ", serviceID)
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return service, err
	}
	return service, nil
}
func GetServiceFromPubkey(pubkey string) {

}
