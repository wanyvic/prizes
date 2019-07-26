package cmd

import (
	"fmt"
	"time"

	prizestypes "github.com/wanyvic/prizes/api/types"
)

//根据时间结算订单
func (s *prizestypes.ServiceOrder) Statement(desiredStatementAt time.Time) error {
	if desiredStatementAt.Before(LastStatementTime) {
		return fmt.Errorf("time invalid Statement time %s LastStatementTime %s", timeAt, s.LastStatementTime)
	}
	statementAt = s.RemoveAt
	if statementAt > desiredStatementAt{
		statementAt = desiredStatementAt	// 订单 定义的删除时间 大于期望删除时间 则不是最后一次结算
	}
	//获取 task 信息
	serviceStatistics, err := ServiceState(s.ServiceID)
	var taskListStatistics []prizestypes.TaskStatistics
	for taskInfo := range serviceStatistics.TaskList {
		if taskInfo.CreatedAt < 
	}
}
