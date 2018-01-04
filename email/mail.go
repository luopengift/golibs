package email

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type Email struct {
	Host     string	`yaml:"host"`
	Username string	`yaml:"username"`
	Password string	`yaml:"password"`
}

func NewEmail(host, username, password string) *Email {
	return &Email{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func (e *Email) auth(mechs string) (smtp.Auth, error) {
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "CRAM-MD5":
			return smtp.CRAMMD5Auth(e.Username, e.Password), nil
		case "PLAIN":
			host, _, err := net.SplitHostPort(e.Host)
			if err != nil {
				return nil, fmt.Errorf("host error:", e.Host)
			}
			return smtp.PlainAuth("", e.Username, e.Password, host), nil
		case "LOGIN":
			return LoginAuth(e.Username, e.Password), nil
		}
	}
	return nil, nil
}

func (e *Email) Send(content *Content) error {
	c, err := smtp.Dial(e.Host)
	if err != nil {
		return err
	}
	defer c.Close()
	if ok, mech := c.Extension("AUTH"); ok {
		auth, err := e.auth(mech)
		if err != nil {
			return err
		}
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("%T failed: %s", auth, err)
			}
		}
	}
	if content.From == "" {
		content.From = e.Username
	}
	if err = c.Mail(content.From); err != nil {
		return err
	}
	for _, to := range strings.Split(content.To, ",") {
		if err = c.Rcpt(to); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = w.Write(content.Bytes()); err != nil {
		return err
	}
	if err = c.Quit(); err.Error() != "250 Ok: queued as " {
		return err
	}
	return nil
}
