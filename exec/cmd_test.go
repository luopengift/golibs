package exec

import (
	"testing"
)

func Test_cmd(t *testing.T) {
	res, err := CmdOut("/bin/bash", "-c", "echo $PATH")
	t.Logf("%v, %v", res, err)
}
