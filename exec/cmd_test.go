package exec

import (
	"github.com/luopengift/golibs/logger"
	"testing"
)

func Test_cmd(t *testing.T) {
	res, err := CmdOut("/bin/bash", "-c", "echo $PATH")
	logger.Info("%v, %v", res, err)
}
