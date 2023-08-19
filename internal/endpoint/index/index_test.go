package index

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	array := make([]int, 0)
	for i := 0; i < 3; i++ {
		array = append(array, i)
	}
	fmt.Println(array)
	// 从0开始切1个.
	temp := array[0:1]
	fmt.Println(temp)
}
