package golibs

import (
    "fmt"
    "github.com/astaxie/beego/httplib"
)

type ApiProvider struct {
    Url string
    Key string
}

func SendSMS(server ApiProvider, phone string, msg string) string {
    req := httplib.Post(server.Url)
    req.Param("apikey", server.Key)
    req.Param("mobile", phone)
    req.Param("text", msg)

    str, err := req.String()
    if err != nil {
        fmt.Println(str)
    }
    return str
}
