package sys

import (
    "os/user"
    "fmt"
    "bytes"
    "github.com/luopengift/golibs/exec"
)

func User() *user.User {
    if user, err := user.Current(); err != nil {
        fmt.Println(err)
        return nil
    } else {
        return user
    }
}

//获取系统的当前用户
func Username() string {
    return User().Username
}

//获取当前用户的家目录
func HomeDir() string {
    return User().HomeDir
}

//检测用户是否存在
//Param: user<string>:用户
//return: bool
func UserExist(user string) bool {
    res, err := exec.CmdOut(fmt.Sprintf("grep ^%s: /etc/passwd | awk -F: '{print $1}'|head -1",user))
    if err != nil {
        return false
    }
    return bool(string(bytes.Trim(res,"\n")) == user)
}
