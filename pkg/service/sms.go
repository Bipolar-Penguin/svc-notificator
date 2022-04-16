package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/go-kit/log"
)

const (
	smsAeroAPIAuthURL string = "https://%s@gate.smsaero.ru/v2/auth"
	smsAeroAPISendURL string = "https://%s@gate.smsaero.ru/v2/sms/send"
)

type smsNotificator struct {
	client     *http.Client
	authString string
	logger     log.Logger
}

func NewSMSNotificator(authString string, logger log.Logger) *smsNotificator {
	s := &smsNotificator{
		authString: authString,
		client:     new(http.Client),
		logger:     logger,
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(smsAeroAPIAuthURL, authString), nil)
	if err != nil {
		s.logger.Log("error", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		s.logger.Log("error", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.logger.Log("event", "not 200 status-code", "status", res.StatusCode)
	}

	return s
}

func (s smsNotificator) NotifyUser(event *domain.Event) {
	defer s.logger.Log("event", "sms done")
	var resBody = map[string]string{"number": "79967726643", "text": translateEvent(event), "sign": "SMS Aero"}

	jsonBody, err := json.Marshal(resBody)
	if err != nil {
		s.logger.Log("error", err)
		return
	}

	url := fmt.Sprintf(smsAeroAPISendURL, s.authString)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.logger.Log("error", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	//	res, err := s.client.Do(req)
	//	if err != nil {
	//		s.logger.Log("error", err)
	//		return
	//	}
	//
	//	response, err := ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		s.logger.Log("error", err)
	//		return
	//	}

	//s.logger.Log("sent", string(response))
}
