package util

import (
	"fmt"
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
