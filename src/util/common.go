package util

import (
	"async/src/constant"
	"fmt"
	"time"
)

func RecoverAsError() error {
	if e := recover(); e != nil {
		var err error
		switch x := e.(type) {
		case string:
			err = fmt.Errorf(x)
		case error:
			err = x
		default:
			err = fmt.Errorf("unknow panic")
		}

		return err
	}
	return nil
}

func GetTimeoutTimer(timeout *time.Duration) *time.Timer {
	var timer *time.Timer
	if timeout == nil {
		timer = time.NewTimer(time.Duration(constant.DefaultTimeoutInSeconds) * time.Second)
	} else {
		timer = time.NewTimer(*timeout)
	}

	return timer
}
