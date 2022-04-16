package trading_bid

import (
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

//func (m *mongoRepository) FindByUserID(userID string) ([]domain.TradingBid, error) {
//	var tradingBids []domain.TradingBid
//
//	//m.coll.Find
//}
