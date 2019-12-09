package util

import (
	"gopkg.in/gomail.v2"
)

var client gomail.SendCloser

func Dial(host string, port int, username, password string) {
	d := gomail.NewDialer(host, port, username, password)
	if c, err := d.Dial(); err != nil {
		Panic("create mail client failed. %v", err)
	} else {
		client = c
	}
}

func SendMail(from, to, cc, bcc string, subject, bodyType, bodyString string) {
	if client == nil {
		Error("missing mail client")
	} else {
		var rec []string
		if ValidMail(from) == false {
			Error("missing mail's sender")
			return
		}
		if ValidMail(to) == false {
			Error("missing mail's to")
			return
		}
		rec = append(rec, to)
		//
		if len(cc) > 0 {
			if ValidMail(cc) == false {
				Error("unknown mail's cc. %s", cc)
				return
			}
			rec = append(rec, to)
		}

		//
		if len(bcc) > 0 {
			if ValidMail(bcc) == false {
				Error("unknown mail's bcc. %s", bcc)
				return
			}
			rec = append(rec, bcc)
		}
		//
		m := gomail.NewMessage()
		m.SetHeader("From", from)
		m.SetHeader("To", to)
		if len(cc) > 0 {
			m.SetAddressHeader("Cc", cc, cc)
		}
		if len(bcc) > 0 {
			m.SetAddressHeader("Bcc", bcc, bcc)
		}
		m.SetHeader("Subject", subject)
		m.SetBody(bodyType, bodyString)
		// Send the email to Bob, Cora and Dan.
		if err := client.Send(from, rec, m); err != nil {
			Error("send mail failed, %v", err)
		}
	}
}

func DirectSendMail(host string, port int, username, password string,
	from, to, cc, bcc string, subject, bodyType, bodyString string) {
	if ValidMail(from) == false {
		Error("missing mail's sender")
		return
	}
	if ValidMail(to) == false {
		Error("missing mail's to")
		return
	}
	if len(cc) > 0 && ValidMail(cc) == false {
		Error("unknown mail's cc. %s", cc)
		return
	}
	if len(bcc) > 0 && ValidMail(bcc) == false {
		Error("unknown mail's bcc. %s", bcc)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetAddressHeader("Cc", cc, cc)
	}
	if len(bcc) > 0 {
		m.SetAddressHeader("Bcc", bcc, bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, bodyString)
	//m.Attach("/home/Alex/lolcat.jpg")
	d := gomail.NewDialer(host, port, username, password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		Error("send mail failed, %v", err)
	}
}
