package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rpc/internal/models"
	"rpc/internal/models/deprecated"
)

const (
	dbName             = "progressData"
	dbUsersCollection  = "m3UserData"
	dbLevelsCollection = "m3LevelData"
	URL                = "mongodb://127.0.0.1:27017"
)

type MongoDAO struct {
	u *mongo.Collection
	l *mongo.Collection
}

func NewMongoDao(ctx context.Context) (*MongoDAO, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDAO{
		u: client.Database(dbName).Collection(dbUsersCollection),
		l: client.Database(dbName).Collection(dbLevelsCollection),
	}, nil
}

func (m *MongoDAO) SetNewRecord(level int32, username string, score int32) error {

	filter := bson.M{"Level": level}

	userScore := deprecated.UserScore{
		Username: username,
		Score:    int(score),
	}

	update := bson.M{"$push": bson.M{"Scores": userScore}}

	options := options.Update().SetUpsert(true)

	_, err := m.l.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDAO) GetBestN(count int, level int) ([]*models.UserScore, error) {

	config := []bson.M{
		{
			"$match": bson.M{
				"Level": level,
			},
		},
		{
			"$unwind": "$Scores",
		},
		{
			"$sort": bson.M{
				"Scores.score": -1,
			},
		},

		{"$project": bson.M{
			"Scores.username": 1,
			"Scores.score":    1,
			"_id":             0,
		}},
		{
			"$limit": count,
		}}

	cursor, err := m.l.Aggregate(context.Background(), config)
	if err != nil {
		return nil, err
	}

	var results []*models.UserScore

	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		var userScore = new(models.UserScore)
		scoreMap := document["Scores"].(primitive.M)

		userScore.Score = scoreMap["score"].(int32)

		userScore.Username = scoreMap["username"].(string)

		results = append(results, userScore)
	}

	return results, nil
}

func (m *MongoDAO) GetUser(key string) *deprecated.User {
	filer := bson.M{"username": key}
	var user deprecated.User
	err := m.u.FindOne(context.Background(), filer).Decode(&user)
	if err != nil {
		return nil
	}

	return &deprecated.User{
		Username: user.Username,
		Password: user.Password,
	}
}
func (m *MongoDAO) CreateUser(username string, password string) error {

	user := deprecated.User{Username: username, Password: password}

	userDoc, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	_, err = m.u.InsertOne(context.Background(), userDoc)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDAO) CheckUser(username string, password string) *deprecated.User {

	filer := bson.M{"username": username,
		"password": password}
	var user deprecated.User
	err := m.u.FindOne(context.Background(), filer).Decode(&user)
	if err != nil {
		return nil
	}

	return &deprecated.User{
		Username: user.Username,
		Password: user.Password,
	}
}
