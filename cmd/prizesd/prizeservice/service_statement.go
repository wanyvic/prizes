package prizeservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	prizestypes "github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
)

var (
	StatementDuration = time.Duration(5 * time.Minute)
)

func Statement(prizeService *service.PrizesService, serviceStatistics prizestypes.ServiceStatistics, LastCheckTime time.Time, NewCheckTime time.Time) (*order.Statement, service.ServiceState, error) {
	if NewCheckTime.After(time.Now().UTC()) {
		return nil, service.ServiceStateUnknown, errors.New("time is too early to statement")
	}

	for i := 0; i < len(prizeService.Order); i++ {
		if prizeService.Order[i].OrderState == order.OrderStatePaying {

			if len(serviceStatistics.TaskList) <= 0 {
				prizeService.NextCheckTime = prizeService.NextCheckTime.Add(StatementDuration)
				prizeService.LastCheckTime = prizeService.Order[i].LastStatementTime
				if prizeService.NextCheckTime.After(prizeService.Order[i].CreatedAt.Add(prizeService.Order[i].TotalTimeDuration)) {
					prizeService.NextCheckTime = prizeService.Order[i].CreatedAt.Add(prizeService.Order[i].TotalTimeDuration)
				}
				return nil, service.ServiceStateUnknown, errors.New("no task need be statement")
			}
			statement, orderState := statementOrder(&prizeService.Order[i], &serviceStatistics, LastCheckTime, NewCheckTime)
			if orderState == order.OrderStateHasBeenPaid {
				if i == len(prizeService.Order)-1 {
					prizeService.State = service.ServiceStateCompleted
				} else {
					prizeService.Order[i+1].OrderState = order.OrderStatePaying
					prizeService.Order[i+1].LastStatementTime = prizeService.Order[i].LastStatementTime
					prizeService.LastCheckTime = prizeService.Order[i+1].LastStatementTime
					prizeService.NextCheckTime = prizeService.NextCheckTime.Add(StatementDuration)
					if prizeService.NextCheckTime.After(prizeService.Order[i+1].CreatedAt.Add(prizeService.Order[i+1].TotalTimeDuration)) {
						prizeService.NextCheckTime = prizeService.Order[i+1].CreatedAt.Add(prizeService.Order[i+1].TotalTimeDuration)
					}
				}
			} else {
				prizeService.NextCheckTime = prizeService.NextCheckTime.Add(StatementDuration)
				prizeService.LastCheckTime = prizeService.Order[i].LastStatementTime
				if prizeService.NextCheckTime.After(prizeService.Order[i].CreatedAt.Add(prizeService.Order[i].TotalTimeDuration)) {
					prizeService.NextCheckTime = prizeService.Order[i].CreatedAt.Add(prizeService.Order[i].TotalTimeDuration)
				}
			}
			return statement, prizeService.State, nil
		}
	}
	return nil, service.ServiceStateUnknown, errors.New("no order or no order state paying")
}

//Statement order
func statementOrder(serviceOrder *order.ServiceOrder, serviceStatistics *prizestypes.ServiceStatistics, LastCheckTime time.Time, NewCheckTime time.Time) (*order.Statement, order.OrderState) {
	TotalUseTime := time.Duration(0)
	taskStatisticsColation := serviceStatistics.TaskList
	amount := int64(0)
	options := order.StatementOptions{
		MasterNodeFeeRate:    serviceOrder.MasterNodeFeeRate,
		DevFeeRate:           serviceOrder.DevFeeRate,
		MasterNodeFeeAddress: serviceOrder.MasterNodeFeeAddress,
		DevFeeAddress:        serviceOrder.DevFeeAddress,
	}

	// compute actually task running time
	for _, taskStatistic := range taskStatisticsColation {
		TotalUseTime += taskStatistic.EndAt.Sub(LastCheckTime)
	}
	// compute statement amount
	amount = int64(TotalUseTime.Hours() * float64(serviceOrder.ServicePrice))
	if amount > serviceOrder.Balance {
		amount = serviceOrder.Balance
	}
	if TotalUseTime > serviceOrder.RemainingTimeDuration {
		TotalUseTime = serviceOrder.RemainingTimeDuration
	}
	logrus.Debug("statement", serviceOrder.RemainingTimeDuration, TotalUseTime)
	if TotalUseTime == serviceOrder.RemainingTimeDuration {
		amount = serviceOrder.Balance
		serviceOrder.OrderState = order.OrderStateHasBeenPaid
		logrus.Debug("lastest order statement")
	}
	statementInfo := parseStatement(taskStatisticsColation, LastCheckTime, NewCheckTime, amount, &options)

	statementInfo.TotalUseTime = TotalUseTime
	serviceOrder.Statement = append(serviceOrder.Statement, *statementInfo)
	serviceOrder.LastStatementTime = NewCheckTime
	serviceOrder.Balance -= amount
	serviceOrder.RemainingTimeDuration -= TotalUseTime
	return &serviceOrder.Statement[len(serviceOrder.Statement)-1], serviceOrder.OrderState
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
		useTime := time.Duration(taskInfo.EndAt.Sub(taskInfo.StartAt))
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
