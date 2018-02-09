package email

import (
	"bytes"
	//	"io"
	"fmt"
	"mime/multipart"
	//"net/mail"
	"net/textproto"
	"time"
)

type Content struct {
	From        string    `json:"From"`
	To          string    `json:"To"`
	Cc          string    `json:"Cc"`
	Bcc         string    `json:"Bcc"`
	Subject     string    `json:"Subject"`
	ContentType string    `json:"Content-Type"`
	Date        time.Time `json:"Date"`
	Body        string    `json:"Body"`
}

func NewContent() *Content {
	return &Content{
		Date: time.Now(),
	}
}

func (ctn *Content) SetFrom(user string) {
	ctn.From = user
}

func (ctn *Content) SetTo(users string) {
	ctn.To = users
}

func (ctn *Content) SetSubject(subject string) {
	ctn.Subject = subject
}

func (ctn *Content) SetContentType(ct string) {
	ctn.ContentType = ct
}

func (ctn *Content) SetBody(body string) error {
	buffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(buffer)
	defer multipartWriter.Close()
	ctn.SetContentType("multipart/alternative;  boundary=" + multipartWriter.Boundary() + "\r\n")
	w, err := multipartWriter.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/html; charset=UTF-8"}})
	if err != nil {
		return fmt.Errorf("creating part for text template: %s", err)
	}
	if _, err = w.Write([]byte(body)); err != nil {
		return err
	}
	ctn.Body = buffer.String()
	return nil
}

func (ctn Content) String() string {
	return string(ctn.Bytes())
}

func (ctn *Content) Bytes() []byte {
	var content bytes.Buffer
	fmt.Fprintf(&content, "From: %s\r\n", ctn.From)
	fmt.Fprintf(&content, "To: %s\r\n", ctn.To)
	fmt.Fprintf(&content, "Subject: %s\r\n", ctn.Subject)
	fmt.Fprintf(&content, "Content-Type: %s\r\n", ctn.ContentType)
	fmt.Fprintf(&content, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&content, "Date: %s\r\n", ctn.Date.Format(time.RFC1123Z))

	fmt.Fprintf(&content, "\r\n")

	fmt.Fprintf(&content, "%v", ctn.Body)

	return content.Bytes()
}
