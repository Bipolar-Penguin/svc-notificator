package service

import (
	"fmt"
	"net/http"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/go-kit/log"
)

const (
	telegramAPIURL string = `https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s`
)

type telegramNotificator struct {
	client   *http.Client
	logger   log.Logger
	botToken string
}

func NewTelegramNotificator(botToken string, logger log.Logger) *telegramNotificator {
	return &telegramNotificator{
		client:   new(http.Client),
		botToken: botToken,
		logger:   logger,
	}
}

func (tn telegramNotificator) NotifyUser(event *domain.Event, telegramID string) {
	url := fmt.Sprintf(telegramAPIURL, tn.botToken, telegramID, translateEvent(event))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tn.logger.Log("error", err)
		return
	}

	_, err = tn.client.Do(req)
	if err != nil {
		tn.logger.Log("error", err)
	}
}
