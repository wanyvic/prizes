package cmd

import (
	"testing"
)

func Test_task_timeuse(t *testing.T) {
	_, err := ServiceState("toqjci8q0jsh46ieh8exp1e3o")
	if err != nil {
		t.Error(err)
	}

}
