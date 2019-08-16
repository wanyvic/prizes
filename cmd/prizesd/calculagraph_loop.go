package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/prizes/cmd"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
)

func CheckCalculagraph() error {
	sign := refresh.NewSign()
	for {
		// logrus.Debug("CheckCalculagraph loop")
		if len(calculagraph.PrioritySequence) > 0 {
			if calculagraph.PrioritySequence[len(calculagraph.PrioritySequence)-1].CheckAt.Before(time.Now().UTC()) {
				item := calculagraph.Pop()
				logrus.Info("CheckCalculagraph check one ", item.ServiceID, " ", item.CheckAt)
				if _, err := cmd.ServiceStatement(item.ServiceID, item.CheckAt); err != nil {
					return err
				}
				continue
			}
			if wait(sign) {
				logrus.Info("CheckCalculagraph exit")
				return nil
			}
			continue
		}
		if wait(sign) {
			logrus.Info("CheckCalculagraph exit")
			return nil
		}
	}
}
func wait(sign *refresh.Signal) bool {
	for i := refresh.TimeScale; i > 0; {
		if i-time.Second > 0 {
			time.Sleep(time.Second)
			i -= time.Second
		} else {
			time.Sleep(i)
			i = 0
		}
		if sign.CheckSign() {
			return true
		}
	}
	return false
}
