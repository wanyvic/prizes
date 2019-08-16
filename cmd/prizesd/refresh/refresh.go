package refresh

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/sirupsen/logrus"
	prizeservice "github.com/wanyvic/prizes/api/types/service"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
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
func (r *RefreshMoudle) Start() {
	go r.whileLoop()
}
func (r *RefreshMoudle) whileLoop() error {
	sign := NewSign()
	for {
		// logrus.Debug("Refreshing docker data to database")
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
				logrus.Info("refresh exit")
				return nil
			}
		}
	}
}
func RefreshStopService(serviceID string) error {
	// logrus.Debug("RefreshStopService")
	taskList, err := db.DBimplement.FindTaskList(serviceID)
	if err != nil {
		return err
	}
	for _, task := range *taskList {
		if task.DesiredState == swarm.TaskStateRunning {
			task.Status.Timestamp = time.Now()
			task.DesiredState = swarm.TaskStateShutdown
			if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RefreshMoudle) refreshDockerTaskFromService(serviceID string) error {
	// logrus.Debug("refreshDockerTaskFromService")
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
		// logrus.Debug("\ttask: ", task.ID)
		if _, err := db.DBimplement.UpdateTaskOne(task); err != nil {
			return err
		}
	}
	return nil
}
func (r *RefreshMoudle) refreshDockerService() error {
	// logrus.Debug("refreshDockerService")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}

	for _, service := range servicelist {
		validNameFilter := filters.NewArgs()
		validNameFilter.Add("service", service.ID)
		tasklist, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: validNameFilter})
		if err != nil {
			return err
		}

		serviceTime, err := db.DBimplement.FindStateTimeAxisOne(service.ID)
		if err != nil {
			logrus.Warning(err)
			if strings.Contains(err.Error(), "no documents in result") {
				serviceTime = &prizeservice.ServiceTimeLine{ServiceID: service.ID}
			} else {
				return err
			}
		}
		serviceTime = updateTimeAxis(serviceTime, tasklist)
		// logrus.Debug("\tservice: ", service.ID)
		if _, err := db.DBimplement.UpdateServiceOne(service); err != nil {
			return err
		}
		if _, err := db.DBimplement.UpdateStateTimeAxisOne(*serviceTime); err != nil {
			return err
		}
		if err := r.refreshDockerTaskFromService(service.ID); err != nil {
			return err
		}
	}
	return nil
}
func (r *RefreshMoudle) refreshDockerNode() error {
	// logrus.Debug("refreshDockerNode")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	nodelist, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodelist {
		// logrus.Debug("\tnode: ", node.ID)
		if _, err := db.DBimplement.UpdateNodeOne(node); err != nil {
			return err
		}
	}
	return nil
}
func updateTimeAxis(serviceTime *prizeservice.ServiceTimeLine, tasklist []swarm.Task) *prizeservice.ServiceTimeLine {
	for _, task := range tasklist {
		if task.DesiredState == swarm.TaskStateRunning {
			if len(serviceTime.TimeAxis) > 0 && serviceTime.TimeAxis[len(serviceTime.TimeAxis)-1].EndAt.Before(time.Unix(0, 0).UTC()) {
				lastAxisgo := &serviceTime.TimeAxis[len(serviceTime.TimeAxis)-1]
				if lastAxisgo.TaskID != task.ID || lastAxisgo.StatusState != task.Status.State {
					nowTime := time.Now().UTC()
					lastAxisgo.EndAt = nowTime
				} else {
					// logrus.Debug("updateTimeAxis same")
					continue
				}
			}
			if task.Status.State == swarm.TaskStateRunning {
				logrus.Debug("updateTimeAxis new asis")
				axis := prizeservice.StateTimeAxis{
					TaskID:       task.ID,
					Version:      task.Meta.Version.Index,
					NodeID:       task.NodeID,
					StartAt:      time.Now().UTC(),
					DesiredState: task.DesiredState,
					StatusState:  task.Status.State,
					Msg:          task.Status.Message,
					Err:          task.Status.Err,
				}
				serviceTime.TimeAxis = append(serviceTime.TimeAxis, axis)
			}
			break
		}
	}
	return serviceTime
}
