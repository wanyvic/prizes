package cmd

import (
	"testing"
)

func Test_task_timeuse(t *testing.T) {
	_, err := ServiceTimeUsed("p31jbl95wm5uhfx9hk79k1w68")
	if err != nil {
		t.Error(err)
	}

}
