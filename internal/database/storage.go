package storage

import (
	"rpc/internal/models"
	"rpc/internal/models/deprecated"
)

type Storage interface {
	GetUser(key string) *deprecated.User
	CreateUser(username string, password string) error
	CheckUser(username string, password string) *deprecated.User

	SetNewRecord(level int32, username string, score int32) error
	GetBestN(count int, level int) ([]*models.UserScore, error)
}
