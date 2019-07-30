package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/cmd/db"
)

//	根据ServiceID 去数据库中获取 所有task 状态、用时、地址等信息
//
func ServiceState(serviceID string, statementAt time.Time) (types.ServiceStatistics, error) {
	logrus.Debug("ServiceState: ", serviceID)
	var serviceStatistics types.ServiceStatistics
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return serviceStatistics, err
	}
	taskList, err := db.DBimplement.FindTaskList(service.ID)
	if err != nil {
		return serviceStatistics, err
	}
	serviceStatistics.ServiceID = serviceID
	serviceStatistics.CreatedAt = service.Meta.CreatedAt
	var taskStatistics []types.TaskStatistics

	for _, task := range *taskList {
		removeTime := task.Status.Timestamp
		if task.DesiredState != swarm.TaskStateShutdown {
			removeTime = statementAt
			serviceStatistics.State = swarm.TaskStateRunning
		}
		var strAddr string
		p_Addr, _ := getAddress(task.NodeID)
		if p_Addr != nil {
			strAddr = *p_Addr
		}
		taskStatistics = append(taskStatistics,
			types.TaskStatistics{
				TaskID:         task.ID,
				NodeID:         task.NodeID,
				CreatedAt:      task.Meta.CreatedAt,
				RemoveAt:       removeTime,
				ReceiveAddress: strAddr,
				State:          task.Status.State,
				Msg:            task.Status.Message,
				Err:            task.Status.Err,
				DesiredState:   task.DesiredState,
			})
		logrus.Debug(fmt.Sprintf("%s %s %s %s %s %s %s ", task.ID, task.NodeID, task.Meta.CreatedAt, removeTime, removeTime.Sub(task.CreatedAt), task.DesiredState, strAddr))
	}
	logrus.Info(fmt.Sprintf("serviceStatistics %s CreatedAt %s usetime %s state %s", serviceID, serviceStatistics.CreatedAt, statementAt.Sub(serviceStatistics.CreatedAt), serviceStatistics.State))
	serviceStatistics.TaskList = taskStatistics
	return serviceStatistics, nil
}
func getAddress(nodeID string) (*string, error) {
	if nodeID == "" {
		return nil, errors.New("empty nodeID")
	}
	node, err := GetNodeInfo(nodeID)
	if err != nil {
		return nil, err
	}
	if value, ok := node.Description.Engine.Labels["REVENUE_ADDRESS"]; ok {
		return &value, nil
	}
	return nil, errors.New("Address not found")
}
