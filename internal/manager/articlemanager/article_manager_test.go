package articlemanager

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	str := "2021-01-11 01:14:53.377"
	// 2021年1月11日 01:14
	res, err := time.ParseInLocation(time.DateTime, str, time.Local)
	if err != nil {
		t.Errorf("ERR:[%+v]", err)
	}
	fmt.Println(res)
}
