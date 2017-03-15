package tail

import (
	"fmt"
	"testing"
)

func Test_tail(t *testing.T) {
	tt := NewTail("test.log")
	tt.ReadLine()
	for v := range tt.NextLine() {
		fmt.Println(*v)
	}
	tt.Stop()
}
