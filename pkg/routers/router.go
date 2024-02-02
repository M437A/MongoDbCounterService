package routers

import (
	"MongoDBCounterService/internal/controllers"
	"MongoDBCounterService/internal/middleware"
	"MongoDBCounterService/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateRouters(router *chi.Mux) {
	newCollectionsRepository := repository.NewCollectionsRepository()
	newCollectionMiddleware := middleware.NewCollectionMiddleware(newCollectionsRepository)
	newCollectionController := controllers.NewCollectionController(newCollectionMiddleware)

	router.Route("/", func(r chi.Router) {
		r.Get("/", newCollectionController.GetNumberOfFiles)
		r.Get("/database", newCollectionController.GetNumberOfFilesByDatabaseName)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}
