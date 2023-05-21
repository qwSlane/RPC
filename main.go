package main

import (
	"context"
	"main/api"
	"main/storage"
)

func main() {

	ctx := context.Background()

	dao, err := storage.NewMongoDao(ctx)
	if err != nil {
		return
	}

	middleware := api.NewAuthMiddleware(dao)
	server := api.NewServer(":8080", dao, middleware)

	server.Start()

}
