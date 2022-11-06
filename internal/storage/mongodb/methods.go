package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"go-cloud-camp/internal/common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateConfig function
func (mb *MongoBackend) CreateConfig(data *common.RequestData) error {
	collFilter := bson.D{{Key: "name", Value: data.Service}}

	// Проверка, существует ли конфиг в базе данных
	collList, err := mb.mdb.ListCollectionNames(context.Background(), collFilter)
	if err != nil {
		return err
	}

	// Если да, отправляем ошибку
	if len(collList) > 0 {
		return common.ErrAlreadyCreated
	}

	if err := mb.mdb.CreateCollection(context.Background(), data.Service); err != nil {
		return err
	}

	coll := mb.mdb.Collection(data.Service)

	counter := bson.D{
		{Key: "_id", Value: "version_counter"},
		{Key: "count", Value: 2},
	}
	if _, err := coll.InsertOne(context.Background(), counter); err != nil {
		return err
	}

	newConfig := &ConfigDataModel{
		CreatedAt: time.Now(),
		ReadedAt:  time.Now(),
		Data:      data.Data,
		Version:   1,
	}

	if _, err := coll.InsertOne(context.Background(), newConfig); err != nil {
		return err
	}

	return nil
}

// ReadConfig function
func (mb *MongoBackend) ReadConfig(service string, version int) ([]byte, error) {
	// Если номер версии больше нуля, тогда добавляем в фильтр поиска
	filter := primitive.D{}
	if version > 0 {
		filter = bson.D{{Key: "version", Value: version}}
	}

	// Обновляем время последнего обращения к конфигу
	update := bson.D{{Key: "$currentDate", Value: bson.D{
		{Key: "readedAt", Value: true},
	}}}

	// Порядок сортировки по убыванию номера версии
	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "version", Value: -1}})

	coll := mb.mdb.Collection(service)
	result := coll.FindOneAndUpdate(context.Background(), filter, update, opts)
	if result.Err() != nil && errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, common.ErrNotFound
	}

	b := &ConfigDataModel{}

	err := result.Decode(b)
	if err != nil {
		return nil, err
	}

	return b.Data, nil
}

// UpdateConfig function
func (mb *MongoBackend) UpdateConfig(data *common.RequestData) error {
	// TODO
	// Здесь используется два обращения к базе данных:
	// инкремент счетчика версий и сохранение нового конфига в базу.
	//
	// НУЖНО обернуть их в транзакцию.
	// НО транзакции в mongodb работают только если сервер в режиме ReplicaSet.
	// В Standalone режиме пока используем последовательные запросы...

	// В работе это может проявиться следующим образом:
	// Если не удалось сохранить новую версию конфига в базу,
	// счетчик версий все равно будет последовательно увеличиваться и
	// некоторые номера версий просто будут пропущены

	if !json.Valid(data.Data) {
		return common.ErrNotValidJsonData
	}

	filterCounter := bson.D{{Key: "_id", Value: "version_counter"}}
	updateCounter := bson.D{{Key: "$inc", Value: bson.D{{
		Key: "count", Value: 1,
	}}}}

	resultCounter := mb.mdb.Collection(data.Service).FindOneAndUpdate(context.Background(), filterCounter, updateCounter)
	if resultCounter.Err() != nil {
		return resultCounter.Err()
	}

	version := &CounterModel{}
	if err := resultCounter.Decode(version); err != nil {
		return err
	}

	newConfig := &ConfigDataModel{
		Version:   version.Count,
		Data:      data.Data,
		CreatedAt: time.Now(),
	}

	if _, err := mb.mdb.Collection(data.Service).InsertOne(context.Background(), newConfig); err != nil {
		return err
	}

	return nil
}

// DeleteConfig function
func (mb *MongoBackend) DeleteConfig(service string, version int) error {
	filter := primitive.D{}
	if version > 0 {
		filter = bson.D{{Key: "version", Value: version}}
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "readedAt", Value: -1}})

	coll := mb.mdb.Collection(service)

	findResult := coll.FindOne(context.Background(), filter, nil, opts)
	if findResult.Err() != nil {
		if findResult.Err() == mongo.ErrNoDocuments {
			return coll.Drop(context.Background())
		} else {
			return findResult.Err()
		}
	}

	configData := &ConfigDataModel{}
	if findErr := findResult.Decode(configData); findErr != nil {
		return findErr
	}

	if time.Since(configData.ReadedAt) < 10*time.Second {
		return common.ErrConfigIsUsed
	}

	var deleteErr error
	if version == 0 {
		deleteErr = coll.Drop(context.Background())
	} else {
		_, deleteErr = coll.DeleteOne(context.Background(), filter)
	}
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
