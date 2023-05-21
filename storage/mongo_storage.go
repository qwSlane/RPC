package storage

import (
	"context"
	"main/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName             = "progressData"
	dbUsersCollection  = "m3UserData"
	dbLevelsCollection = "m3LevelData"
	URL                = "mongodb://127.0.0.1:27017"
)

type mongoDAO struct {
	u *mongo.Collection
	l *mongo.Collection
}

func NewMongoDao(ctx context.Context) (*mongoDAO, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mongoDAO{
		u: client.Database(dbName).Collection(dbUsersCollection),
		l: client.Database(dbName).Collection(dbLevelsCollection),
	}, nil
}

func (m *mongoDAO) SetNewRecord(level int, username string, score int) error {

	filter := bson.M{"Level": level}

	userScore := types.UserScore{
		Username: username,
		Score:    score,
	}

	update := bson.M{"$push": bson.M{"Scores": userScore}}

	options := options.Update().SetUpsert(true)

	_, err := m.l.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDAO) GetBestN(count int, level int) ([]types.UserScore, error) {

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

	var results []types.UserScore

	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		var userScore types.UserScore
		scoreMap := document["Scores"].(primitive.M)

		userScore.Username = scoreMap["username"].(string)
		userScore.Score = int(scoreMap["score"].(int32))

		results = append(results, userScore)
	}

	return results, nil
}

func (m *mongoDAO) GetUser(key string) *types.User {
	filer := bson.M{"username": key}
	var user types.User
	err := m.u.FindOne(context.Background(), filer).Decode(&user)
	if err != nil {
		return nil
	}

	return &types.User{
		Username: user.Username,
		Password: user.Password,
	}
}
func (m *mongoDAO) CreateUser(username string, password string) error {

	user := types.User{Username: username, Password: password}

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

func (m *mongoDAO) CheckUser(username string, password string) *types.User {

	filer := bson.M{"username": username,
		"password": password}
	var user types.User
	err := m.u.FindOne(context.Background(), filer).Decode(&user)
	if err != nil {
		return nil
	}

	return &types.User{
		Username: user.Username,
		Password: user.Password,
	}
}
