package massgrid

import (
	"testing"

	"github.com/wanyvic/prizes/api/types/order"
)

func Test_getblockhash(t *testing.T) {

	DefaultNetParams = DefaultTestNetParams
	refund := order.RefundPayment{}
	refund.Payments = append(refund.Payments, order.Payment{ReceiveAddress: "mob5XbfAsJrSbBVXe9ZiCSy2aMQ7rKqMip", Amount: 13288629856})
	refund.Payments = append(refund.Payments, order.Payment{ReceiveAddress: "mob5XbfAsJrSbBVXe9ZiCSy2aMQ7rKqMip", Amount: 16880000000})
	refund.Payments = append(refund.Payments, order.Payment{ReceiveAddress: "mob5XbfAsJrSbBVXe9ZiCSy2aMQ7rKqMip", Amount: 16880000000})
	refund.Payments = append(refund.Payments, order.Payment{ReceiveAddress: "mob5XbfAsJrSbBVXe9ZiCSy2aMQ7rKqMip", Amount: 16880000000})
	amount := int64(0)
	for _, payment := range refund.Payments {
		amount += payment.Amount
	}
	t.Log(amount)
	s, err := SendMany(&refund)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(s)
	}
}
