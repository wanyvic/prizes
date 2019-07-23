package server

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func Test_Unix_Server(t *testing.T) {
	u, err := NewServer("unix")
	if err != nil {
		t.Error(err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = u.Stop()
	}()
	if err := u.Start(); err != nil {
		t.Error(err)
	}
}
func Test_Tcp_Server(t *testing.T) {
	u, err := NewServer("tcp")
	if err != nil {
		t.Error(err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = u.Stop()
	}()
	if err := u.Start(); err != nil {
		t.Error(err)
	}
}
