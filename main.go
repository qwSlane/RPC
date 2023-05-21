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

	server := api.NewServer(":8080", dao, nil)

	server.Start()

}
