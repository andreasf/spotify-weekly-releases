package platform

import "time"

//go:generate counterfeiter . Time
type Time interface {
	Sleep(d time.Duration)
}

type TimeWrapper struct{}

func (self *TimeWrapper) Sleep(d time.Duration) {
	time.Sleep(d)
}
