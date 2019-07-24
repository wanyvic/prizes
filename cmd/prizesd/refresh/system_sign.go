package refresh

import (
	"os"
	"os/signal"
	"syscall"
)

type Signal struct {
	sign bool
}

func NewSign() *Signal {
	s := &Signal{sign: false}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.sign = true
	}()
	return s
}
func (s *Signal) CheckSign() bool {
	if s.sign {
		return true
	}
	return false
}
