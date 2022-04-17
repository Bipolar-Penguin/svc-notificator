package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/go-kit/log"
)

const (
	smtpPort = 587
	smtpHost = "smtp.yandex.ru"
)

type mailNotificator struct {
	auth     smtp.Auth
	username string
	logger   log.Logger
}

func NewEmailNotificator(username, password string, logger log.Logger) *mailNotificator {
	return &mailNotificator{
		auth:     smtp.PlainAuth("", username, password, smtpHost),
		username: username,
		logger:   logger,
	}
}

func (e mailNotificator) NotifyUser(event *domain.Event, userEmail string) {
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Обновление котировочной сессии\r\n\r\n"+
		"Котировочная сессия обновлена: %s\r\n",
		e.username,
		userEmail,
		TranslateEvent(event),
	)

	clientConn, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		e.logger.Log("error", err)
	}

	err = clientConn.StartTLS(&tls.Config{ServerName: smtpHost})
	if err != nil {
		e.logger.Log("error", err)
		return
	}
	err = clientConn.Auth(e.auth)
	if err != nil {
		e.logger.Log("error", err)
		return
	}

	err = clientConn.Mail(e.username)
	if err != nil {
		e.logger.Log("error", err)
		return
	}

	err = clientConn.Rcpt("kuwerin@gmail.com")
	if err != nil {
		e.logger.Log("error", err)
		return
	}

	wc, err := clientConn.Data()
	if err != nil {
		e.logger.Log("error", err)
		return
	}
	_, err = fmt.Fprintf(wc, msg)
	if err != nil {
		e.logger.Log("error", err)
		return
	}
	err = wc.Close()
	if err != nil {
		e.logger.Log("error", err)
		return
	}

	err = clientConn.Close()
	if err != nil {
		e.logger.Log("error", err)
		return
	}
}
