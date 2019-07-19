package cmd

import (
	"fmt"
	"time"

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
}

func ServiceTimeUsed(serviceID string) (ServiceStatistics, error) {
	var serviceStatistics ServiceStatistics
	service, err := db.DBimplement.FindServiceOne(serviceID)
	if err != nil {
		return serviceStatistics, err
	}
	fmt.Printf("service: %s\n", service.ID)
	taskList, err := db.DBimplement.FindTaskList(service.ID)
	if err != nil {
		return serviceStatistics, err
	}
	serviceStatistics.ServiceID = serviceID
	serviceStatistics.CreatedAt = service.CreatedAt
	var taskStatistics []TaskStatistics
	var td time.Duration
	latestTime := time.Unix(0, 0).UTC()
	fmt.Println("taskID", "nodeID", "CreatedAt", "RemoveAt", "useTime", "state")
	for _, task := range *taskList {
		removeTime := task.Status.Timestamp
		if task.DesiredState != "shutdown" {
			removeTime = time.Unix(0, 0).UTC()
		}
		if task.Status.Timestamp.After(latestTime) {
			latestTime = task.Status.Timestamp
		}
		taskStatistics = append(taskStatistics, TaskStatistics{TaskID: task.ID, NodeID: task.NodeID, CreatedAt: task.Meta.CreatedAt, RemoveAt: removeTime})
		fmt.Println(task.ID, task.NodeID, task.Meta.CreatedAt, removeTime, task.Status.Timestamp.Sub(task.CreatedAt), task.DesiredState)
		td += task.Status.Timestamp.Sub(task.CreatedAt)
	}
	serviceStatistics.TaskList = taskStatistics
	fmt.Printf("service total time %s\n", td)
	return serviceStatistics, nil
}
