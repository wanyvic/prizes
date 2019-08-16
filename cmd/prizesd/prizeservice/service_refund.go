package prizeservice

import (
	"strconv"
	"time"

	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
)

func Refund(prizeService *service.PrizesService) *order.RefundPayment {
	refundPayment := order.RefundPayment{}
	refundPayment.RefundID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + service.DefaultServiceRefundID + CreateRandomNumberString(8)
	refundPayment.CreatedAt = time.Now().UTC()
	var amount int64
	for i := 0; i < len(prizeService.Order); i++ {
		payment := order.Payment{}
		payment.ReceiveAddress = prizeService.Order[i].Drawee
		payment.Amount = prizeService.Order[i].Balance
		refundPayment.Payments = append(refundPayment.Payments, payment)
		amount += prizeService.Order[i].Balance
		prizeService.Order[i].Balance = 0
		prizeService.Order[i].RemainingTimeDuration = 0
	}
	refundPayment.TotalAmount = amount
	prizeService.State = service.ServiceStateCompleted
	return &refundPayment
}
