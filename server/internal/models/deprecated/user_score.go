package deprecated

type UserScore struct {
	Username string `json:"username" bson:"username"`
	Score    int    `json:"score" bson:"score"`
}
