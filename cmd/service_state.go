package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
)

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

func ServiceState(serviceID string) (ServiceStatistics, error) {
	logrus.Info("ServiceState: ", serviceID)
	var serviceStatistics ServiceStatistics
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return serviceStatistics, err
	}
	taskList, err := db.DBimplement.FindTaskList(service.ID)
	if err != nil {
		return serviceStatistics, err
	}
	serviceStatistics.ServiceID = serviceID
	serviceStatistics.CreatedAt = service.CreatedAt
	var taskStatistics []TaskStatistics
	var td time.Duration
	latestTime := time.Unix(0, 0).UTC()

	logrus.Info("taskID", "nodeID", "CreatedAt", "RemoveAt", "useTime", "state", "Address")
	for _, task := range *taskList {
		removeTime := task.Status.Timestamp
		if task.DesiredState != "shutdown" {
			removeTime = time.Unix(0, 0).UTC()

			serviceStatistics.State = "running"
		}
		if task.Status.Timestamp.After(latestTime) {
			latestTime = task.Status.Timestamp
		}
		var strAddr string
		p_Addr, _ := getAddress(task.NodeID)
		if p_Addr != nil {
			strAddr = *p_Addr
		}
		taskStatistics = append(taskStatistics,
			TaskStatistics{
				TaskID:         task.ID,
				NodeID:         task.NodeID,
				CreatedAt:      task.Meta.CreatedAt,
				RemoveAt:       removeTime,
				ReceiveAddress: strAddr,
				State:          string(task.Status.State),
				Msg:            task.Status.Message,
				Err:            task.Status.Err,
				DesiredState:   string(task.DesiredState),
			})
		fmt.Println(task.ID, task.NodeID, task.Meta.CreatedAt, removeTime, task.Status.Timestamp.Sub(task.CreatedAt), task.DesiredState, strAddr)
		td += task.Status.Timestamp.Sub(task.CreatedAt)
	}
	serviceStatistics.TaskList = taskStatistics

	logrus.Info("ervice total time: ", td)
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
