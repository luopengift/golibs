package set

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	set1 := NewSet(1, "2", "3")
	set1.Add(1)
	fmt.Println(set1.Contains(1))
	set2 := NewSet("1", "2", "3", "4")
	fmt.Println(set1, set2)
	set3 := NewSet()
	fmt.Println("2-1", set2.Diff(set1))
	fmt.Println("1-2", set1.Diff(set2))
	fmt.Println(set1.Diff(set3))
	fmt.Println(set3.Len())
}
