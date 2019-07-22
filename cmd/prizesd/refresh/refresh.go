package refresh

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

const (
	DefaultTimeScale = 3000
)

var (
	TimeScale time.Duration
)

func init() {
	TimeScale = time.Duration(DefaultTimeScale) * time.Millisecond
}
func RefreshStopService(serviceID string) error {
	logrus.Info("RefreshStopService")
	taskList, err := db.DBimplement.FindTaskList(serviceID)
	if err != nil {
		return err
	}
	for _, task := range *taskList {
		if task.DesiredState == "running" {
			task.Status.Timestamp = time.Now()
			task.DesiredState = "shutdown"
			if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
				return err
			}
		}
	}
	return nil
}

func refreshDockerTaskFromService(serviceID string) error {
	logrus.Info("refreshDockerTaskFromService")
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("service", serviceID)
	tasklist, err := dockerapi.CLI.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
	if err != nil {
		return err
	}
	for _, task := range tasklist {
		logrus.Info("\ttask: ", task.ID)
		if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
			return err
		}
	}
	return nil
}
func WhileLoop() error {
	logrus.Info("WhileLoop")
	for {
		if err := refreshDockerNode(); err != nil {
			return err
		}
		if err := refreshDockerService(); err != nil {
			return err
		}
		for i := TimeScale; i > 0; {
			if i-time.Second > 0 {
				time.Sleep(time.Second)
				i -= time.Second
			} else {
				time.Sleep(i)
				i = 0
			}
			if dockerapi.FExit {
				return nil
			}
		}
	}
}
func refreshDockerService() error {
	logrus.Info("refreshDockerService")
	servicelist, err := dockerapi.CLI.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}
	for _, service := range servicelist {
		logrus.Info("\tservice: ", service.ID)
		if _, err := db.DBimplement.UpdateServiceOne(service); err != nil {
			return err
		}
		if err := refreshDockerTaskFromService(service.ID); err != nil {
			return err
		}
	}
	return nil
}
func refreshDockerNode() error {
	logrus.Info("refreshDockerNode")
	nodelist, err := dockerapi.CLI.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodelist {
		logrus.Info("\tnode: ", node.ID)
		if _, err := db.DBimplement.UpdateNodeOne(node); err != nil {
			return err
		}
	}
	return nil
}
