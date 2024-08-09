package main

import (
	"components/myconfig"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

type EmailHelper struct {
}

func (this *EmailHelper) SendEmail(subject string, body string) error {
	f, err1 := os.OpenFile("email.html", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err1 != nil {
		fmt.Println(err1)
	}
	f.Write([]byte(body))
	defer f.Close()

	smtpHost := myconfig.GConfig.Smtp.SmtpHost
	smtpPort := myconfig.GConfig.Smtp.SmtpPort
	senderEmail := myconfig.GConfig.Smtp.SenderEmail
	password := myconfig.GConfig.Smtp.SenderPassword

	subject = fmt.Sprintf("Subject: %v\n", subject)
	from := fmt.Sprintf("From: %v\n", senderEmail)
	to := fmt.Sprintf("To: %v\n", myconfig.GConfig.Smtp.Receivers)
	contentType := "Content-Type: text/html; charset=UTF-8\n\n"
	body = fmt.Sprintf(`<html><body>%v</body></html>`, body)

	message := []byte(from + to + subject + contentType + body)

	serverAddress := smtpHost + ":" + smtpPort
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", serverAddress, tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(senderEmail); err != nil {
		return err
	}
	if err := client.Rcpt(myconfig.GConfig.Smtp.Receivers); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(message)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
