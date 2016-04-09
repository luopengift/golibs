# golibs
some basic libs for golang

useage

####coding example
package main

import (
  "github.com/luopengift/golibs"
)

func main() {
    smtp := SMTPServer{Server: "smtp.exmail.qq.com:25", Sender: "****@qq.com", Passwd: "******"}
    mail := Mail{To: "****@qq.com", Type: "html", Subject: "邮件测试", Body: "HELLO"}
    fmt.Println(SendMail(smtp, mail))
    
