package email

import (
	"testing"
	"fmt"
)

func Test_Send(t *testing.T) {
	user := "luopeng@qq.com"
	password := "xxxxxxxxxx"
	host := "smtp.exmail.qq.com:25"
	body := `<html><body><h3>"Test send to email~~"</h3></body></html>`
	content := NewContent()
	content.SetFrom("luopeng@qq.com")
	content.SetTo("870148195@qq.com")
	content.SetSubject("使用Golang发送邮件")
	content.SetBody(body)

	fmt.Println(content)

	mail := NewEmail(host, user, password)
	fmt.Println(mail.Send(content))
}
