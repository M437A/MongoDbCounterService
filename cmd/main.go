package main

import (
	"MongoDBCounterService/pkg/connections"
	"MongoDBCounterService/pkg/routers"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	connections.Init(ctx)
	defer connections.CloseMongoDBConnection()

	router := chi.NewRouter()
	routers.CreateRouters(router)

	err := http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
