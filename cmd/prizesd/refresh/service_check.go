package refresh

import (
	"container/heap"
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
)

var (
	Loop ServiceLoop
	var lck sync.Mutex
)

type CheckItem struct {
	ServiceID string `json:"service_id,omitempty"`
	RemoveAt  time.Time
}
type ServiceLoop []CheckItem

func (h ServiceLoop) Len() int { return len(h) }

func (h ServiceLoop) Less(i, j int) bool { return h[i].RemoveAt.Before(h[j].RemoveAt) }
func (h ServiceLoop) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ServiceLoop) Push(x interface{}) {
	item := x.(*CheckItem)
	*h = append(*h, *item)
}

func (h *ServiceLoop) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func InitLoop() error {
	cli, err := dockerapi.GetDockerClient()
	if err != nil {
		return err
	}
	servicelist, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return err
	}

	heap.Init(&Loop)
	for _, service := range servicelist {
		removeAt, _ := time.Parse("2006-01-02 15:04:05", service.Spec.Labels["com.massgrid.deletetime"])
		item := CheckItem{ServiceID: service.ID, RemoveAt: removeAt}
		heap.Push(&Loop, &item)
	}
	return nil
}
func ChangeServiceRemoveTime(ServiceID string, newTime time.Time) {
	lck.Lock()
	defer lck.Unlock()
	logrus.Info("ChangeServiceRemoveTime ", ServiceID, " ", newTime)
	for i := 0; i < len(Loop); i++ {
		if Loop[i].ServiceID == ServiceID {
			Loop[i].RemoveAt = newTime
			heap.Fix(&Loop, i)
			return
		}
	}
}
func CheckLoop() error {
	sign := NewSign()
	for {
		logrus.Info("CheckLoop")
		if len(Loop) > 0 {
			if Loop[len(Loop)-1].RemoveAt.After(time.Now().UTC()) {
				lck.Lock()
				item := heap.Pop(&Loop).(*CheckItem)
				lck.Unlock()
				if err := cmd.ServiceRemove(item.ServiceID); err != nil {
					return err
				}
			}
		} else {
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
}
