package refresh

import (
	"context"
	"time"

	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
)

const (
	DefaultTimeScale = 3000
)

var (
	TimeScale time.Duration
)

type RefreshMoudle struct {
	TimeScale time.Duration
}

func init() {
	TimeScale = time.Duration(DefaultTimeScale) * time.Millisecond
}
func NewRefreshMoudle() *RefreshMoudle {
	r := &RefreshMoudle{TimeScale: TimeScale}
	return r
}
func (r *RefreshMoudle) WhileLoop() error {
	sign := NewSign()
	for {
		logrus.Info("Refreshing docker data to database")
		if err := r.refreshDockerNode(); err != nil {
			logrus.Error(err.Error())
		}
		if err := r.refreshDockerService(); err != nil {
			logrus.Error(err.Error())
		}
		for i := TimeScale; i > 0; {
			if i-time.Second > 0 {
				time.Sleep(time.Second)
				i -= time.Second
			} else {
				time.Sleep(i)
				i = 0
			}
			if sign.CheckSign() {
				return nil
			}
		}
	}
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

func (r *RefreshMoudle) refreshDockerTaskFromService(serviceID string) error {
	logrus.Debug("refreshDockerTaskFromService")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("service", serviceID)
	tasklist, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
	if err != nil {
		return err
	}
	for _, task := range tasklist {
		logrus.Debug("\ttask: ", task.ID)
		if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
			return err
		}
	}
	return nil
}
func (r *RefreshMoudle) refreshDockerService() error {
	logrus.Debug("refreshDockerService")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}
	for _, service := range servicelist {
		logrus.Debug("\tservice: ", service.ID)
		if _, err := db.DBimplement.UpdateServiceOne(service); err != nil {
			return err
		}
		if err := r.refreshDockerTaskFromService(service.ID); err != nil {
			return err
		}
	}
	return nil
}
func (r *RefreshMoudle) refreshDockerNode() error {
	logrus.Debug("refreshDockerNode")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	nodelist, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodelist {
		logrus.Debug("\tnode: ", node.ID)
		if _, err := db.DBimplement.UpdateNodeOne(node); err != nil {
			return err
		}
	}
	return nil
}
