package prizeservice

import (
	"strconv"
	"time"

	"github.com/wanyvic/prizes/api/types/order"
	"github.com/wanyvic/prizes/api/types/service"
)

func Refund(prizeService *service.PrizesService) *[]order.RefundPayment {
	refundPaymentArray := []order.RefundPayment{}
	for i := 0; i < len(prizeService.Order); i++ {
		refundPaymentArray = append(refundPaymentArray, refundOrder(&prizeService.Order[i]))
		prizeService.Order[i].OrderState = order.OrderStateHasBeenRefund
	}
	prizeService.State = service.ServiceStateCompleted
	return &refundPaymentArray
}
func refundOrder(serviceOrder *order.ServiceOrder) order.RefundPayment {
	refundPayment := order.RefundPayment{}
	refundPayment.RefundID = strconv.FormatInt(time.Now().UTC().Unix(), 10) + service.DefaultServiceRefundID + CreateRandomNumberString(8)
	refundPayment.CreatedAt = time.Now().UTC()
	refundPayment.TotalAmount = serviceOrder.Balance
	refundPayment.Drawee = serviceOrder.Drawee
	serviceOrder.Refund = &refundPayment
	serviceOrder.Balance = 0
	return refundPayment
}
