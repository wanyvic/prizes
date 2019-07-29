package prizeservice

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/swarm"
	prizestypes "github.com/wanyvic/prizes/api/types"
	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
)

func Statement(prizeService *service.PrizesService, serviceStatistics prizestypes.ServiceStatistics, desiredTime time.Time, options order.StatementOptions) (*order.Statement, error) {
	var statement *order.Statement
	for i := 0; i < len(prizeService.Order); i++ {
		if prizeService.Order[i].OrderState == order.OrderStatePaying {
			statement = statementOrder(&prizeService.Order[i], &serviceStatistics, desiredTime, &options)
			break
		}
	}
	return statement, nil
}
func statementOrder(serviceOrder *order.ServiceOrder, serviceStatistics *prizestypes.ServiceStatistics, desiredTime time.Time, options *order.StatementOptions) *order.Statement {
	taskStatisticsColation := []prizestypes.TaskStatistics{}
	statementAt := desiredTime
	for _, taskStatistics := range serviceStatistics.TaskList {
		if taskStatistics.State == swarm.TaskStateRunning {
			taskStatisticsColation = append(taskStatisticsColation, taskStatistics)
		} else if taskStatistics.RemoveAt.After(serviceOrder.LastStatementTime) {
			taskStatisticsColation = append(taskStatisticsColation, taskStatistics)
		}
	}
	balanceUsableTime := time.Duration(float64(serviceOrder.Balance)/float64(serviceOrder.ServicePrice)) * time.Hour
	if balanceUsableTime < desiredTime.Sub(serviceOrder.LastStatementTime) { //不够结算
		statementAt = serviceOrder.LastStatementTime.Add(balanceUsableTime)
		serviceOrder.OrderState = order.OrderStateHasBeenPaid
	}
	statementInfo := parseStatement(taskStatisticsColation, serviceOrder.LastStatementTime, statementAt, serviceOrder.Balance, options)
	serviceOrder.Statement = append(serviceOrder.Statement, *statementInfo)
	serviceOrder.LastStatementTime = statementAt
	serviceOrder.NextStatementTime = statementAt.Add(options.StatementDuration)
	return statementInfo
}
func parseStatement(taskStatisticsColation []prizestypes.TaskStatistics, statementStartAt time.Time, statementEndAt time.Time, amount int64, options *order.StatementOptions) *order.Statement {

	statement := order.Statement{}
	statement.StatementID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + CreateRandomNumberString(6)
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
		taskAmount := int64(float64(amount) * useTime.Hours() / TotalUseTime.Hours())
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
func CreateRandomNumberString(len int) string {
	var container string
	var str = "1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
