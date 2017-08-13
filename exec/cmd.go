package exec

import (
    "github.com/luopengift/golibs/logger"
    "errors"
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

func Cmd(name string,arg ...string) *exec.Cmd {
    return exec.Command(name, arg...)
}

func CmdOut(name string, arg ...string) (string, error) {
    cmd := exec.Command(name, arg...)
    ret,err := cmd.CombinedOutput()
    return string(ret), err
}


func CmdOutBytes(name string, arg ...string) ([]byte, error) {
    cmd := exec.Command(name, arg...)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return out.Bytes(), err
}

func CmdOutNoLn(name string, arg ...string) (out string, err error) {
    out, err = CmdOut(name, arg...)
    if err != nil {
        return
    }

    return strings.TrimSpace(string(out)), nil
}


func CmdRunWithTimeout(command string, timeout int64) ([]byte, error) {
    cmd := Cmd("/bin/bash","-c",command)
    done := make(chan error)
    var stdout,stderr bytes.Buffer
    cmd.Stdout, cmd.Stderr = &stdout, &stderr

    cmd.Start()
    go func() {
        done <- cmd.Wait()
    }()

    select {
    case <-time.After(time.Duration(timeout) * time.Second):

        go func() {
            result := <-done // allow goroutine to exit
            logger.Info("%v",result)
        }()

        err := cmd.Process.Kill()   // timeout
        if err != nil {
            logger.Error("failed to kill: %s, error: %s", cmd.Path, err)
            return stderr.Bytes(),err
        }
        return stderr.Bytes(),
            errors.New(fmt.Sprintf("exit status 62:TIMEOUT %d,Process %s has been killed",timeout,strings.Join(cmd.Args, " ")))
    case err := <-done:
        //fmt.Println(string(stdout.Bytes()),string(stderr.Bytes()),err)
        if err != nil {
            return stderr.Bytes(), err
        }
        return stdout.Bytes(), err
    }
}
