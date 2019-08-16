package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/cmd/db"
)

//ServiceState returns ServiceStatistics error
func ServiceState(serviceID string, startTime time.Time, endTime time.Time) (types.ServiceStatistics, error) {
	logrus.Debug("ServiceState: ", serviceID)
	var serviceStatistics types.ServiceStatistics
	serviceTimeAxis, err := db.DBimplement.FindStateTimeAxisOne(serviceID)
	if err != nil {
		return serviceStatistics, err
	}
	serviceStatistics.ServiceID = serviceID
	serviceStatistics.StartAt = startTime
	serviceStatistics.EndAt = endTime
	var taskStatistics []types.TaskStatistics
	for _, timeAxis := range serviceTimeAxis.TimeAxis {
		if timeAxis.EndAt.Before(time.Unix(0, 0).UTC()) {
			timeAxis.EndAt = endTime
		}
		if timeAxis.StartAt.Before(startTime) {
			timeAxis.StartAt = startTime
		}
		if timeAxis.EndAt.After(startTime) {
			var strAddr string
			pAddr, _ := getAddress(timeAxis.NodeID)
			if pAddr != nil {
				strAddr = *pAddr
			}
			taskStatistics = append(taskStatistics,
				types.TaskStatistics{
					TaskID:         timeAxis.TaskID,
					NodeID:         timeAxis.NodeID,
					StartAt:        timeAxis.StartAt,
					EndAt:          timeAxis.EndAt,
					ReceiveAddress: strAddr,
					State:          timeAxis.StatusState,
					Msg:            timeAxis.Msg,
					Err:            timeAxis.Err,
					DesiredState:   timeAxis.DesiredState,
				})
			logrus.Debug(fmt.Sprintf("%s %s %s %s %s %s %s ", timeAxis.TaskID, timeAxis.NodeID, timeAxis.StartAt, timeAxis.EndAt, timeAxis.EndAt.Sub(timeAxis.StartAt), timeAxis.DesiredState, strAddr))
		}
	}
	logrus.Info(fmt.Sprintf("serviceStatistics %s CreatedAt %s usetime %s", serviceID, serviceStatistics.StartAt, endTime.Sub(serviceStatistics.EndAt)))
	serviceStatistics.TaskList = taskStatistics
	return serviceStatistics, nil
}

//getAddress get node receive_address
func getAddress(nodeID string) (*string, error) {
	if nodeID == "" {
		return nil, errors.New("empty nodeID")
	}
	node, err := GetNodeInfo(nodeID)
	if err != nil {
		return nil, err
	}
	if value, ok := node.Description.Engine.Labels[types.LabelRevenueAddress]; ok {
		return &value, nil
	}
	return nil, errors.New("Address not found")
}
