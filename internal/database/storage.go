package storage

import (
	"rpc/internal/models/deprecated"
)

type Storage interface {
	GetUser(key string) *deprecated.User
	CreateUser(username string, password string) error
	CheckUser(username string, password string) *deprecated.User

	SetNewRecord(level int, username string, score int) error
	GetBestN(count int, level int) ([]deprecated.UserScore, error)
}
