package file

import (
	"fmt"
	"testing"
)

func Test_read(t *testing.T) {
	//tt := NewTail("test-%Y-%M-%D-%h-%m.log")
	tt := NewTail("test.log", &TimeRuler{})
	tt.ReadLine()
	tt.EndStop(true)
	for v := range tt.NextLine() {
		fmt.Println(string(v))
	}
}

func Test_tail(t *testing.T) {
    f := NewTail("/tmp/test/t_%Y%M%D%h%m.log",&TimeRuler{})
    f.ReadLine()
    f.EndStop(false)
    for v := range f.NextLine() {
        fmt.Println(string(v))
    }
}
