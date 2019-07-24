package server

import (
	"testing"
)

func Test_Unix_Server(t *testing.T) {
	_, err := NewServer("unix")
	if err != nil {
		t.Error(err)
	}

}
func Test_Tcp_Server(t *testing.T) {
	_, err := NewServer("tcp")
	if err != nil {
		t.Error(err)
	}

}
