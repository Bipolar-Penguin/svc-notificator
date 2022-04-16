package user

import "github.com/Bipolar-Penguin/svc-notificator/pkg/domain"

type Repository interface {
	Find(userID string) (domain.User, error)
}
