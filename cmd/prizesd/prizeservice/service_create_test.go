package prizeservice

import (
	"testing"

	"github.com/wanyvic/prizes/api/types/prizeservice"
)

func Test_service_create(t *testing.T) {
	serviceCreate := prizeservice.ServiceCreate{}
	serviceCreate.Amount = 10
	serviceCreate.ServicePrice = 1
	_, err := Create(&serviceCreate)
	if err != nil {
		t.Error(err)
	}
}
