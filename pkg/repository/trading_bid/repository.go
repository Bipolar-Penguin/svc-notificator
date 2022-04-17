package trading_bid

import "github.com/Bipolar-Penguin/svc-notificator/pkg/domain"

type Repository interface {
	FindMany(tradingSessionID string) ([]domain.TradingBid, error)
}
