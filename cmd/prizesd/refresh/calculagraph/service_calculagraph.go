package calculagraph

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd/db"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

var (
	PrioritySequence ServiceCalculagraph
	Lck              sync.Mutex
)

type CheckItem struct {
	ServiceID string `json:"service_id,omitempty"`
	CheckAt   time.Time
}
type ServiceCalculagraph []CheckItem

func (h ServiceCalculagraph) Len() int { return len(h) }

func (h ServiceCalculagraph) Less(i, j int) bool { return h[i].CheckAt.Before(h[j].CheckAt) }
func (h ServiceCalculagraph) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ServiceCalculagraph) Push(x interface{}) {
	item := x.(*CheckItem)
	*h = append(*h, *item)
}

func (h *ServiceCalculagraph) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func InitCalculagraph() error {
	logrus.Debug("InitCalculagraph")
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}
	heap.Init(&PrioritySequence)
	for _, service := range servicelist {
		if prizeService, err := db.DBimplement.FindPrizesServiceOne(service.ID); err == nil {
			if prizeService.NextCheckTime.After(time.Unix(0, 0).UTC()) {
				item := CheckItem{ServiceID: service.ID, CheckAt: prizeService.NextCheckTime}
				heap.Push(&PrioritySequence, &item)
				logrus.Info("InitCalculagraph push one ", service.ID, " ", item.CheckAt)
			}
		}
	}
	return nil
}
func ChangeCheckTime(ServiceID string, newTime time.Time) {
	Lck.Lock()
	defer Lck.Unlock()
	logrus.Info("ChangeCheckTime ", ServiceID, " ", newTime)
	for i := 0; i < len(PrioritySequence); i++ {
		if PrioritySequence[i].ServiceID == ServiceID {
			PrioritySequence[i].CheckAt = newTime
			heap.Fix(&PrioritySequence, i)
			return
		}
	}
	logrus.Info("calculagraph Push ", ServiceID, " ", newTime)
	heap.Push(&PrioritySequence, &CheckItem{ServiceID, newTime})
}
func RemoveService(ServiceID string) {
	Lck.Lock()
	defer Lck.Unlock()
	logrus.Info("RemoveService ", ServiceID)
	for i := 0; i < len(PrioritySequence); i++ {
		if PrioritySequence[i].ServiceID == ServiceID {
			if i == len(PrioritySequence)-1 {
				PrioritySequence = append(PrioritySequence[:i])
			} else {
				PrioritySequence = append(PrioritySequence[:i], PrioritySequence[i+1])
			}
			heap.Fix(&PrioritySequence, i)
			return
		}
	}
}
func Push(ServiceID string, newTime time.Time) {
	Lck.Lock()
	defer Lck.Unlock()
	logrus.Info("calculagraph Push ", ServiceID, " ", newTime)
	heap.Push(&PrioritySequence, &CheckItem{ServiceID, newTime})
}
func Pop() *CheckItem {
	Lck.Lock()
	defer Lck.Unlock()
	item := heap.Pop(&PrioritySequence).(CheckItem)
	return &item
}
