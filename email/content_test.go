package email

import (
	"fmt"
	"testing"
)

func Test_Content(t *testing.T) {
	content := Content{
		From: "from@example.com",
		To:   "To@example",
	}
	fmt.Println(content)
}
