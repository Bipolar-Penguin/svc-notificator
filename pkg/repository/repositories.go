package repository

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/repository/user"
)

const (
	tradingDatabase string = "trading"

	userCollection string = "users"
)

type repositories struct {
	User user.Repository
}

func MakeRepositories(mongoURL string, logger log.Logger) (*repositories, error) {
	var r = new(repositories)

	var err error

	clientOpts := options.Client().ApplyURI(mongoURL)

	clientOpts.SetServerSelectionTimeout(30 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		logger.Log("error", err)
		return nil, err
	}

	r.User = user.NewMongoRepository(client.Database(tradingDatabase).Collection(userCollection))
	r.User = user.LoggingMiddleware(logger)(r.User)

	return r, nil
}
