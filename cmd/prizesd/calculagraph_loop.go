package main

import (
	"container/heap"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
)

func CheckCalculagraph() error {
	sign := NewSign()
	for {
		logrus.Info("CheckCalculagraph")
		if len(calculagraph.PrioritySequence) > 0 {
			if calculagraph.PrioritySequence[len(calculagraph.PrioritySequence)-1].RemoveAt.After(time.Now().UTC()) {
				Lck.Lock()
				item := heap.Pop(&calculagraph.PrioritySequence).(*CheckItem)
				Lck.Unlock()
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
