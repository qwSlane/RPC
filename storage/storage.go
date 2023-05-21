package storage

import "main/types"

type Storage interface {
	GetUser(key string) *types.User
	CreateUser(username string, password string) error
	CheckUser(username string, password string) *types.User

	SetNewRecord(level int, username string, score int) error
	GetBestN(count int, level int) ([]types.UserScore, error)
}
