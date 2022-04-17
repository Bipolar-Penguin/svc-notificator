package service

import (
	"fmt"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/Bipolar-Penguin/svc-notificator/pkg/repository"
	"github.com/go-kit/log"
)

type Notificator interface {
	NotifyUser(event *domain.Event, userID string)
}

func TranslateEvent(event *domain.Event) string {
	return fmt.Sprintf("Ваша ставка неактуальна новая минимальная ставка: %d", event.Amount)
}

type service struct {
	telegram Notificator
	sms      Notificator
	mail     Notificator
	rep      *repository.Repositories
	logger   log.Logger
}

func NewService(telegram telegramNotificator, sms smsNotificator, mail mailNotificator, rep *repository.Repositories, logger log.Logger) *service {
	return &service{
		telegram: telegram,
		sms:      sms,
		mail:     mail,
		rep:      rep,
		logger:   logger,
	}
}

func (s service) NotifyUser(event *domain.Event, _ string) {
	s.logger.Log("event", fmt.Sprintf("%v", event), "location", "notificator")

	tradingBids, err := s.rep.TradingBid.FindMany(event.EventID)
	if err != nil {
		s.logger.Log("error", err)
		return
	}

	var userIDs = make(map[string]int)
	for _, bid := range tradingBids {
		if bid.UserID != event.GUID {
			userIDs[bid.UserID] = 0
		}
	}

	for userID, _ := range userIDs {
		user, err := s.rep.User.Find(userID)
		if err != nil {
			s.logger.Log("error", err)
			break
		}

		fmt.Printf("%v", user)
		if s.telegram != nil && user.Permissions.Telegram && user.Contacts.TelegramID != "" {
			s.logger.Log("event", "notifying via telegram")
			s.telegram.NotifyUser(event, user.Contacts.TelegramID)
		}

		if s.sms != nil && user.Permissions.Phone && user.Contacts.PhoneNumber != "" {
			s.logger.Log("event", "notifying via sms")
			s.sms.NotifyUser(event, user.Contacts.PhoneNumber)
		}

		if s.mail != nil && user.Permissions.Email && user.Contacts.Email != "" {
			s.logger.Log("event", "notifying via email")
			s.mail.NotifyUser(event, user.Contacts.Email)
		}
	}

}
