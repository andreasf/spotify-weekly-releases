package platform

import "time"

//go:generate counterfeiter . Time
type Time interface {
	Sleep(d time.Duration)
	Now() time.Time
}

type TimeWrapper struct{}

func (self *TimeWrapper) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (self *TimeWrapper) Now() time.Time {
	return time.Now()
}