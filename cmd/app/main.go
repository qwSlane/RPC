package main

import (
	"context"
	"rpc/internal/app"
	"rpc/internal/database"
	"rpc/internal/services/middleware"
)

func main() {

	ctx := context.Background()

	dao, err := storage.NewMongoDao(ctx)
	if err != nil {
		return
	}

	mw := middleware.NewAuthMiddleware(dao)
	server := app.NewServer(":8080", dao, mw)

	err = server.Start()
	if err != nil {
		return
	}

}
