package user

import (
	"context"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(coll *mongo.Collection) *mongoRepository {
	return &mongoRepository{
		coll: coll,
	}
}

func (m *mongoRepository) Find(userID string) (domain.User, error) {
	var user domain.User

	var err error

	err = m.coll.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, nil
	}
	if err != nil {
		return user, err
	}

	return user, nil
}
