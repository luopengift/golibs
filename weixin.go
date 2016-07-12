package golibs

import (
    "fmt"
    "log"
    "strings"
    "github.com/astaxie/beego/httplib"
    "crypto/tls"
    "time"
    "encoding/json"
    "net/url"
)

type WeixinApiProvider struct {
    Url string
    //Key string
    Appid string
    Appsecret string
}

func SendWeiXin(server WeixinApiProvider, wxname string, msg string) string {
    log.Println("Send SMS now")
    fmt.Println("SendSMS print by geothe",wxname,msg,server)	

    // req := httplib.Post(server.Url)
    //req.Param("apikey", server.Key)
    req := httplib.Get(server.Url + "appid=" + server.Appid + "&appsecret=" + server.Appsecret + "")
    req.SetTimeout(5*time.Second, 5*time.Second)
    req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
    //req.Param("appid",server.Appid)
    //req.Param("appsecret",server.Appsecret)
    //req.Param("verify","false")

    reqstr,err :=req.String()
    if err !=nil {
	fmt.Println(reqstr)
    }

    result := make(map[string]string)
    json.Unmarshal([]byte(reqstr), &result)
    access_token, _ := result["access_token"] 

    fmt.Println("accesstoken is----------->>>>>",access_token)

    wxurl :="https://10.205.140.41:11443/v1/wechatqy/send/?"
    wxq :=httplib.Get(wxurl + "username="+url.QueryEscape(strings.Trim(strings.Replace(wxname,",","|",-1),",")) + "&content=" + url.QueryEscape(msg) +"")
    wxq.SetTimeout(5*time.Second, 5*time.Second)
    wxq.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
    wxq.Header("Authorization", access_token)
    //wxq.Param("username", strings.Trim(strings.Replace(wxname,",","|",-1),","))   //change , to |
    //wxq.Param("content", msg)

    fmt.Println("luopengift printlog",wxq)

    str, err := wxq.String()
  fmt.Println("luopengift printlog",str)
    if err != nil {
        fmt.Println(str)
    }
    return str
}
