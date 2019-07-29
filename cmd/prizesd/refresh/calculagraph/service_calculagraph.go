package calculagraph

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

var (
	PrioritySequence ServiceCalculagraph
	Lck              sync.Mutex
)

type CheckItem struct {
	ServiceID string `json:"service_id,omitempty"`
	RemoveAt  time.Time
}
type ServiceCalculagraph []CheckItem

func (h ServiceCalculagraph) Len() int { return len(h) }

func (h ServiceCalculagraph) Less(i, j int) bool { return h[i].RemoveAt.Before(h[j].RemoveAt) }
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
		removeAt, _ := time.Parse("2006-01-02 15:04:05", service.Spec.Labels["com.massgrid.deletetime"])
		item := CheckItem{ServiceID: service.ID, RemoveAt: removeAt}
		heap.Push(&PrioritySequence, &item)
	}
	return nil
}
func ChangeServiceRemoveTime(ServiceID string, newTime time.Time) {
	Lck.Lock()
	defer Lck.Unlock()
	logrus.Info("ChangeServiceRemoveTime ", ServiceID, " ", newTime)
	for i := 0; i < len(PrioritySequence); i++ {
		if PrioritySequence[i].ServiceID == ServiceID {
			PrioritySequence[i].RemoveAt = newTime
			heap.Fix(&PrioritySequence, i)
			return
		}
	}
}
func Push(ServiceID string, newTime time.Time) {
	Lck.Lock()
	defer Lck.Unlock()
	logrus.Info("Push ", ServiceID, " ", newTime)
	heap.Push(&PrioritySequence, CheckItem{ServiceID, newTime})
}
