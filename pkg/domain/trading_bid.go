package domain

import "time"

type TradingBid struct {
	TradingSessionID string    `json:"trading_session_id"`
	UserID           string    `json:"user_id"`
	Date             time.Time `json:"date"`
	Salary           int       `json:"salary"`
}
