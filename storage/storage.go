package storage

import "main/types"

type Storage interface {
	GetUser(key string) *types.User
	CreateUser(username string, password string) error
	Delete(key string) error
	CheckUser(username string, password string) *types.User
}
