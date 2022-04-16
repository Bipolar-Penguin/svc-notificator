package service

import (
	"fmt"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/go-kit/log"
)

type Notificator interface {
	NotifyUser(event *domain.Event)
}

func translateEvent(event *domain.Event) string {
	return fmt.Sprintf("Ваша ставка неактуальна: %s", event.Action)
}

type service struct {
	telegram Notificator
	sms      Notificator
	mail     Notificator
	logger   log.Logger
}

func NewService(telegram telegramNotificator, sms smsNotificator, mail mailNotificator, logger log.Logger) *service {
	return &service{
		telegram: telegram,
		sms:      sms,
		mail:     mail,
		logger:   logger,
	}
}

func (s service) NotifyUser(event *domain.Event) {
	if s.telegram != nil {
		s.telegram.NotifyUser(event)
	}

	if s.sms != nil {
		s.sms.NotifyUser(event)
	}

	if s.mail != nil {
		s.mail.NotifyUser(event)
	}
}
