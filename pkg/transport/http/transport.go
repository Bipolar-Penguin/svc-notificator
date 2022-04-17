package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/websocket"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
)

const (
	defaultReadBufferSize  int = 1024
	defaultWriteBufferSize int = 1024
)

type websocketServer struct {
	addr           int
	tradingUpdates <-chan *domain.Event
	logger         log.Logger
}

func NewWebsocketServer(addr int, tradingUpdates <-chan *domain.Event, logger log.Logger) *websocketServer {
	return &websocketServer{
		addr:           addr,
		tradingUpdates: tradingUpdates,
		logger:         logger,
	}
}

func (s *websocketServer) StreamUpdatesToWebsocket() {
	http.HandleFunc("/ws/trading", func(rw http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  defaultReadBufferSize,
			WriteBufferSize: defaultWriteBufferSize,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		ws, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			s.logger.Log("error", err)
		}

		for {
			select {

			case <-s.tradingUpdates:
				for event := range s.tradingUpdates {
					msg, err := json.Marshal(event)
					if err != nil {
						s.logger.Log("error", err)
					}

					if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
						s.logger.Log("error", err)
					}
				}
			}

		}

	})

	go http.ListenAndServe(fmt.Sprintf(":%d", s.addr), nil)
}
