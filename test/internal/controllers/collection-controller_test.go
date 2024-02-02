package controllers__test

import (
	"MongoDBCounterService/internal/controllers"
	"MongoDBCounterService/internal/model"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockMiddleware struct{}

func (m *MockMiddleware) GetNumberOfFilesFromAllDB(ctx context.Context) ([]model.CollectionStats, error) {
	return []model.CollectionStats{
		{CollectionName: "collection1", DocumentCount: 30},
		{CollectionName: "collection2", DocumentCount: 20},
	}, nil
}

func (m *MockMiddleware) GetListOfCollectionsStats(database string, ctx context.Context) ([]model.CollectionStats, error) {
	return []model.CollectionStats{
		{CollectionName: "collection1", DocumentCount: 25},
		{CollectionName: "collection2", DocumentCount: 15},
	}, nil
}

func TestGetNumberOfFiles(t *testing.T) {
	controller := controllers.NewCollectionController(&MockMiddleware{})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	controller.GetNumberOfFiles(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `[{"CollectionName":"collection1","DocumentCount":30},{"CollectionName":"collection2","DocumentCount":20}]`
	assert.Equal(t, expectedBody, strings.TrimSuffix(w.Body.String(), "\n"))
}

func TestGetNumberOfFilesByDatabaseName(t *testing.T) {
	controller := controllers.NewCollectionController(&MockMiddleware{})
	requestBody := `{"database":"test_name"}`
	req := httptest.NewRequest("GET", "/database", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	controller.GetNumberOfFilesByDatabaseName(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `[{"CollectionName":"collection1","DocumentCount":25},{"CollectionName":"collection2","DocumentCount":15}]`
	assert.Equal(t, expectedBody, strings.TrimSuffix(w.Body.String(), "\n"))
}
