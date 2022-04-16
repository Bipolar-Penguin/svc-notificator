package user

import (
	"time"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Repository
}

func LoggingMiddleware(logger log.Logger) func(Repository) Repository {
	return func(next Repository) Repository {
		return &loggingMiddleware{logger, next}
	}
}

func (l *loggingMiddleware) Find(userID string) (user domain.User, err error) {
	defer func(begin time.Time) {
		l.logger.Log("entity", "user", "method", "find", "error", err)
	}(time.Now())
	return l.next.Find(userID)
}
