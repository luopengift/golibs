package file

import (
	"fmt"
	"testing"
)

func Test_read(t *testing.T) {
	//tt := NewTail("test-%Y-%M-%D-%h-%m.log")
	tt := NewTail("test.log")
	tt.ReadLine()

	for v := range tt.NextLine() {
		fmt.Println(*v)
	}
	tt.Stop()
}
