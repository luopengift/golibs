package message

import (
        "fmt"
        "net/smtp"
        "strings"
)

func SendToMail(username, password, host, to, subject, body, mailtype string) error {
        auth := smtp.PlainAuth("", username, password, strings.Split(host, ":")[0])
        var contentType string
        if mailtype == "html" {
                contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
        } else {
                contentType = "Content-Type: text/plain" + "; charset=UTF-8"
        }

        msg := []byte("To: " + to + "\r\n" + "From: " + user + ">\r\nSubject: " + "\r\n" + contentType + "\r\n\r\n" + body)
        send_to := strings.Split(to, ";")
        err := smtp.SendMail(host, auth, user, send_to, msg)
        return err
}
