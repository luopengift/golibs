package exec

import (
    "testing"
    "github.com/luopengift/golibs/logger"
)

func Test_cmd(t *testing.T) {
    res,err := CmdOut("/bin/bash","-c","echo $PATH")
    logger.Info("%v, %v", res, err)
}
