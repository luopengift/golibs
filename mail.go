package golibs

import (
    "net/smtp"
    "strings"
)

type SMTPServer struct {
    Server string
    Sender string
    Passwd string
}

type Mail struct {
    To      string
    Type    string
    Subject string
    Body    string
}

func SendMail(server SMTPServer, mail Mail) error {
    auth := smtp.PlainAuth("", server.Sender, server.Passwd, strings.Split(server.Server, ":")[0])
    var Type string
    switch mail.Type {
    case "html":
        Type = "Content-Type: text/" + mail.Type + "; charset=UTF-8"
    default:
        Type = "Content-Type: text/plain" + "; charset=UTF-8"
    }
    msg := []byte("To: " + mail.To + "\r\n" + "From: " + server.Sender + ">\r\n" + "Subject: " + "\r\n" + Type + "\r\n\r\n" + mail.Body)
    err := smtp.SendMail(server.Server, auth, server.Sender, strings.Split(mail.To, ";"), msg)
    return err
}
