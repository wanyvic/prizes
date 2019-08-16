package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/massgrid"
	"github.com/wanyvic/prizes/cmd/prizesd/prizeservice"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
)

//ServiceCreate returns response and error
func ServiceCreate(serviceCreate *service.ServiceCreate) (*types.ServiceCreateResponse, error) {
	prizeService, response, err := prizeservice.Create(serviceCreate)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.Push(prizeService.DockerService.ID, prizeService.NextCheckTime)
	return response, nil
}

//ServiceUpdate returns response and error
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
	calculagraph.ChangeCheckTime(prizeService.DockerService.ID, prizeService.NextCheckTime)
	return response, nil
}

//ServiceStatement returns statement and error
func ServiceStatement(ServiceID string, statementAt time.Time) (*order.Statement, error) {
	prizeService, err := db.DBimplement.FindPrizesServiceOne(ServiceID)
	if err != nil {
		return nil, err
	}
	if prizeService.State == service.ServiceStateCompleted {
		return nil, errors.New("service has been statement")
	}
	serviceStatistics, err := ServiceState(ServiceID, prizeService.LastCheckTime, statementAt)
	if err != nil {
		return nil, err
	}
	statement, serviceState, err := prizeservice.Statement(prizeService, serviceStatistics, prizeService.LastCheckTime, statementAt)
	if err != nil {
		logrus.Warning(err)
	}
	if statement != nil {
		hash, err := massgrid.SendMany(statement)
		if err != nil {
			return nil, err
		}
		statement.StatementTransaction = *hash
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	if serviceState == service.ServiceStateRunning {
		calculagraph.ChangeCheckTime(prizeService.DockerService.ID, prizeService.NextCheckTime)
	} else if serviceState == service.ServiceStateCompleted {
		err := serviceRemove(ServiceID)
		calculagraph.RemoveService(ServiceID)
		if err != nil {
			return nil, err
		}
	}
	return statement, nil
}

//ServiceRefund returns RefundInfo and error
func ServiceRefund(ServiceID string) (*order.RefundPayment, error) {
	var err error
	refundPayment := &order.RefundPayment{}

	_, err = ServiceStatement(ServiceID, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	prizeService, err := db.DBimplement.FindPrizesServiceOne(ServiceID)
	if err != nil {
		return nil, err
	}
	if prizeService.State == service.ServiceStateCompleted {
		return nil, errors.New("service has been paid")
	}
	refundPayment = prizeservice.Refund(prizeService)

	hash, err := massgrid.SendMany(refundPayment)
	if err != nil {
		return nil, err
	}
	refundPayment.RefundTransaction = *hash
	err = serviceRemove(ServiceID)
	if err != nil {
		return nil, err
	}
	_, err = db.DBimplement.UpdatePrizesServiceOne(*prizeService)
	if err != nil {
		return nil, err
	}
	calculagraph.RemoveService(ServiceID)

	logrus.Info(fmt.Sprintf("%+v", refundPayment))
	return refundPayment, nil
}

func ServiceInfo(serviceID string) (*service.ServiceInfo, error) {
	prizeService, err := db.DBimplement.FindPrizesServiceOne(serviceID)
	if err != nil {
		return nil, err
	}
	serviceInfo, err := prizeservice.ServiceInfo(prizeService)
	if err != nil {
		return nil, err
	}
	return serviceInfo, nil
}

func GetServicesFromPubkey(pubkey string) (*[]service.ServiceInfo, error) {
	serviceInfoList := []service.ServiceInfo{}
	prizeServiceList, err := db.DBimplement.FindPrizesServiceFromPubkey(pubkey)
	if err != nil {
		return nil, err
	}
	for _, prizeService := range *prizeServiceList {
		serviceInfo, err := prizeservice.ServiceInfo(&prizeService)
		if err != nil {
			return nil, err
		}
		serviceInfoList = append(serviceInfoList, *serviceInfo)
	}
	return &serviceInfoList, nil
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
