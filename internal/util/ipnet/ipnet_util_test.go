package ipnetutil

import (
	"fmt"
	"testing"
)

func TestParseIPLocation(t *testing.T) {
	res, err := ParseIPLocation("113.90.48.5")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
