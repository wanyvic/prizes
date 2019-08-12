package massgrid

import (
	"testing"
)

func Test_getblockhash(t *testing.T) {
	err := getblockhash(1)
	t.Error(err)
}
