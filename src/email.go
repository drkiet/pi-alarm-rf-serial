package main

import (
	"strings"
	"fmt"
	"io"
	"bufio"
	"os"
	"crypto/tls"
	"log"
	"net/smtp"
)


type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body string
}

type SmtpServer struct {
	Host      string
	Port      string
	TlsConfig *tls.Config
}

var routineEmail string
var routinePsw string
var toList map[string]string

// email code is from here: 
// https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0 
func (s *SmtpServer) ServerName() string {
	return s.Host + ":" + s.Port
}

func (mail *Mail) BuildMessage() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", mail.Sender)
	if len(mail.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}
	if len(mail.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	header += "\r\n" + mail.Body

	return header
}


func sendEmail(to string, subject string, message string) {
	mail := Mail{}
	mail.Sender = routineEmail
	mail.To = []string{to}
	mail.Cc = []string{}
	mail.Bcc = []string{}
	mail.Subject = subject
	mail.Body = message

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{Host: "smtp.gmail.com", Port: "465"}
	smtpServer.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	auth := smtp.PlainAuth("", mail.Sender, routinePsw, smtpServer.Host)
	conn, err := tls.Dial("tcp", smtpServer.ServerName(), smtpServer.TlsConfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.Host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		log.Panic(err)
	}
	receivers := append(mail.To, mail.Cc...)
	receivers = append(receivers, mail.Bcc...)
	for _, k := range receivers {
		log.Println("sending to: ", k)
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}

// Initialize email configuration
func emailInit() {
	toList = make(map[string]string)

	file, err := os.Open(PiConfigFile)
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file for email/text: .private")
	}

	reader := bufio.NewReader(bufio.NewReader(file))

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("ReadLine failed: ", err)
			}
			break
		}
		parseEmailSettings(string(bytes))
	}
	fmt.Println(toList)
}

func getToList() (map[string]string) {
	return toList;
}

func addTo (toCfg string) {
	toTokens := strings.Split(toCfg, "=")
	toList[strings.Trim(toTokens[0], " ")] = strings.Trim(toTokens[1], " ")
}

func parseEmailSettings(line string) {
	tokens := strings.Split(line, ":")

	switch strings.Trim(tokens[0], " ") {
	case "email":
		routineEmail = (strings.Trim(tokens[1], " "))

	case "psw":
		routinePsw = (strings.Trim(tokens[1], " "))

	case "to":
		addTo(strings.Trim(tokens[1], " "))

	default:
	}
}
