package controllers

import (
	"MongoDBCounterService/internal/middleware"
	"MongoDBCounterService/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CollectionController struct {
	middleware middleware.Collection
}

func NewCollectionController(middleware middleware.Collection) *CollectionController {
	return &CollectionController{middleware: middleware}
}

func (c *CollectionController) GetNumberOfFiles(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("Get request for counting collections in MongoDB")
	ctx := context.Background()

	mongoCollections, serviceException := c.middleware.GetNumberOfFilesFromAllDB(ctx)
	if JSONResponseException(serviceException, responseWriter) {
		log.Printf("Service error: %s", serviceException.Error())
		return
	}

	log.Printf("All collections were counted")
	stoutResponse(mongoCollections)
	JSONResponse(responseWriter, mongoCollections)
}

func (c *CollectionController) GetNumberOfFilesByDatabaseName(responseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("Get request for counting collections in MongoDB")
	var databaseName model.DataBaseName

	err := json.NewDecoder(request.Body).Decode(&databaseName)
	if JSONResponseException(err, responseWriter) {
		return
	}

	database := databaseName.DataBaseName
	ctx := context.Background()

	mongoCollections, serviceException := c.middleware.GetListOfCollectionsStats(database, ctx)
	if JSONResponseException(serviceException, responseWriter) {
		log.Printf("Service error: %s", serviceException.Error())
		return
	}

	log.Printf("All collections were counted")
	stoutResponse(mongoCollections)
	JSONResponse(responseWriter, mongoCollections)
}

func JSONResponse(responseWriter http.ResponseWriter, data interface{}) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	err := json.NewEncoder(responseWriter).Encode(data)
	if err != nil {
		log.Printf("Error encoding JSON response: %s", err.Error())
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
	}
}

func JSONResponseException(err error, responseWriter http.ResponseWriter) bool {
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return true
	}
	return false
}

func stoutResponse(slice []model.CollectionStats) {
	if len(slice) == 0 {
		return
	}

	var sb strings.Builder

	for _, stats := range slice {
		sb.WriteString(stats.CollectionName)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%d", stats.DocumentCount))
		sb.WriteString("\n")
	}

	log.Println(sb.String())
}
