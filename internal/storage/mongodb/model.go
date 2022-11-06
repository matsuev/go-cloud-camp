package mongodb

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConfigDataModel struct
type ConfigDataModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Version   int                `bson:"version"`
	CreatedAt time.Time          `bson:"createdAt"`
	ReadedAt  time.Time          `bson:"readedAt"`
	Data      json.RawMessage    `bson:"data"`
}

// CounterModel struct
type CounterModel struct {
	ID    string `bson:"_id"`
	Count int    `bson:"count"`
}
