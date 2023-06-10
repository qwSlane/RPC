package deprecated

type Level struct {
	Level  int         `bson:"level"`
	Scores []UserScore `bson:"scores"`
}
