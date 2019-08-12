package prizeservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/swarm"
	prizestypes "github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
)

var (
	StatementDuration = time.Duration(5 * time.Minute)
)

func Statement(prizeService *service.PrizesService, serviceStatistics prizestypes.ServiceStatistics, desiredTime time.Time) (*order.Statement, service.ServiceState, error) {
	var statement *order.Statement
	for i := 0; i < len(prizeService.Order); i++ {
		if prizeService.Order[i].OrderState == order.OrderStatePaying {
			statement := statementOrder(&prizeService.Order[i], &serviceStatistics, desiredTime)
			if prizeService.Order[i].OrderState == order.OrderStateHasBeenPaid {
				if i == len(prizeService.Order)-1 {
					prizeService.State = service.ServiceStateCompleted
				} else {
					prizeService.Order[i+1].OrderState = order.OrderStatePaying
					prizeService.Order[i+1].LastStatementTime = prizeService.Order[i].LastStatementTime
					prizeService.NextCheckTime = prizeService.NextCheckTime.Add(StatementDuration)
					if prizeService.NextCheckTime.After(prizeService.Order[i+1].RemoveAt) {
						prizeService.NextCheckTime = prizeService.Order[i+1].RemoveAt
					}
				}
			} else {
				prizeService.NextCheckTime = prizeService.NextCheckTime.Add(StatementDuration)
				if prizeService.NextCheckTime.After(prizeService.Order[i].RemoveAt) {
					prizeService.NextCheckTime = prizeService.Order[i].RemoveAt
				}
			}
			return statement, prizeService.State, nil
		}
	}
	return statement, prizeService.State, errors.New("no order or no order state paying")
}
func statementOrder(serviceOrder *order.ServiceOrder, serviceStatistics *prizestypes.ServiceStatistics, desiredTime time.Time) *order.Statement {
	taskStatisticsColation := []prizestypes.TaskStatistics{}
	statementAt := desiredTime
	var statementInfo *order.Statement
	amount := int64(0)
	for _, taskStatistics := range serviceStatistics.TaskList {
		if taskStatistics.State == swarm.TaskStateRunning {
			taskStatisticsColation = append(taskStatisticsColation, taskStatistics)
		} else if taskStatistics.RemoveAt.After(serviceOrder.LastStatementTime) {
			taskStatisticsColation = append(taskStatisticsColation, taskStatistics)
		}
	}
	balanceUsableTime := time.Duration(float64(serviceOrder.Balance) / float64(serviceOrder.ServicePrice) * float64(time.Hour))

	options := order.StatementOptions{
		MasterNodeFeeRate:    serviceOrder.MasterNodeFeeRate,
		DevFeeRate:           serviceOrder.DevFeeRate,
		MasterNodeFeeAddress: serviceOrder.MasterNodeFeeAddress,
		DevFeeAddress:        serviceOrder.DevFeeAddress,
	}
	logrus.Debug("statement", balanceUsableTime, desiredTime.Sub(serviceOrder.LastStatementTime))
	if balanceUsableTime <= desiredTime.Sub(serviceOrder.LastStatementTime)+time.Minute { //不够结算
		statementAt = serviceOrder.LastStatementTime.Add(balanceUsableTime)
		serviceOrder.OrderState = order.OrderStateHasBeenPaid
		amount = serviceOrder.Balance
		statementInfo = parseStatement(taskStatisticsColation, serviceOrder.LastStatementTime, statementAt, amount, &options)
	} else {
		logrus.Debug("not latest statement")
		amount = int64(statementAt.Sub(serviceOrder.LastStatementTime).Hours() * float64(serviceOrder.ServicePrice))
		statementInfo = parseStatement(taskStatisticsColation, serviceOrder.LastStatementTime, statementAt, amount, &options)
	}
	serviceOrder.Statement = append(serviceOrder.Statement, *statementInfo)
	serviceOrder.LastStatementTime = statementAt
	serviceOrder.Balance -= amount
	return statementInfo
}

func parseStatement(taskStatisticsColation []prizestypes.TaskStatistics, statementStartAt time.Time, statementEndAt time.Time, amount int64, options *order.StatementOptions) *order.Statement {

	statement := order.Statement{}
	statement.StatementID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + service.DefaultStatementID + CreateRandomNumberString(8)
	statement.CreatedAt = time.Now().UTC()
	statement.StatementStartAt = statementStartAt
	statement.StatementEndAt = statementEndAt
	statement.TotalAmount = amount
	statement.MasterNodeFeeRate = options.MasterNodeFeeRate
	statement.DevFeeRate = options.DevFeeRate
	statement.MasterNodeFeeAddress = options.MasterNodeFeeAddress
	statement.DevFeeAddress = options.DevFeeAddress
	masterNodeAmount := int64(0)
	DevAmount := int64(0)
	TotalUseTime := time.Duration(statementEndAt.Sub(statementStartAt))

	if statement.MasterNodeFeeRate != 0 {
		masterNodeAmount = int64(amount * statement.MasterNodeFeeRate / 10000.0)
	}
	if statement.DevFeeRate != 0 {
		DevAmount = int64(amount * statement.DevFeeRate / 10000.0)
		DevPayment := order.Payment{ReceiveAddress: statement.DevFeeAddress, Amount: DevAmount, Msg: "Dev pay"}
		statement.Payments = append(statement.Payments, DevPayment)
	}
	amount = amount - masterNodeAmount - DevAmount
	for _, taskInfo := range taskStatisticsColation {
		if taskInfo.CreatedAt.Before(statementStartAt) {
			taskInfo.CreatedAt = statementStartAt
		}
		if taskInfo.State == swarm.TaskStateRunning {
			taskInfo.RemoveAt = statementEndAt
		}
		useTime := time.Duration(taskInfo.RemoveAt.Sub(taskInfo.CreatedAt))
		taskAmount := int64(float64(amount)*useTime.Hours()/TotalUseTime.Hours() + 0.5)
		var msg string
		if taskInfo.ReceiveAddress == "" {
			taskInfo.ReceiveAddress = statement.MasterNodeFeeAddress
			msg = "node receive address not found, replaced by masternode"
		}
		taskPayment := order.Payment{ReceiveAddress: taskInfo.ReceiveAddress, Amount: taskAmount, TaskID: taskInfo.TaskID, TaskState: taskInfo.State, Msg: msg}
		statement.Payments = append(statement.Payments, taskPayment)
	}
	if statement.MasterNodeFeeRate != 0 {
		masterNodePayment := order.Payment{ReceiveAddress: statement.MasterNodeFeeAddress, Amount: masterNodeAmount, Msg: "masternode pay"}
		statement.Payments = append(statement.Payments, masterNodePayment)
	}
	return &statement
}
