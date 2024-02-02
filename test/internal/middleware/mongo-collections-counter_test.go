package middleware_test

import (
	"MongoDBCounterService/internal/middleware"
	"MongoDBCounterService/internal/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCollectionsRepository struct {
	mock.Mock
}

func (m *MockCollectionsRepository) ListDatabaseNames(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCollectionsRepository) ListCollectionNames(databaseName string, ctx context.Context) ([]string, error) {
	args := m.Called(databaseName, ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCollectionsRepository) LengthOfCollection(databaseName, collectionName string, ctx context.Context) (int64, error) {
	args := m.Called(databaseName, collectionName, ctx)
	return args.Get(0).(int64), args.Error(1)
}

// todo: make tests for check exceptions,empty and big slices
func TestGetNumberOfFilesFromAllDB(t *testing.T) {
	repoMock := new(MockCollectionsRepository)

	expectedDatabases := []string{"db1", "db2", "db3"}
	repoMock.On("ListDatabaseNames", mock.Anything).Return(expectedDatabases, nil)

	repoMock.On("ListCollectionNames", "db1", mock.Anything).Return([]string{"collection1", "collection2", "collection3"}, nil)
	repoMock.On("ListCollectionNames", "db2", mock.Anything).Return([]string{"collection4", "collection5", "collection6"}, nil)
	repoMock.On("ListCollectionNames", "db3", mock.Anything).Return([]string{"collection7", "collection8", "collection9"}, nil)

	repoMock.On("LengthOfCollection", "db1", "collection1", mock.Anything).Return(int64(10), nil)
	repoMock.On("LengthOfCollection", "db1", "collection2", mock.Anything).Return(int64(60), nil)
	repoMock.On("LengthOfCollection", "db1", "collection3", mock.Anything).Return(int64(100), nil)

	repoMock.On("LengthOfCollection", "db2", "collection4", mock.Anything).Return(int64(90), nil)
	repoMock.On("LengthOfCollection", "db2", "collection5", mock.Anything).Return(int64(40), nil)
	repoMock.On("LengthOfCollection", "db2", "collection6", mock.Anything).Return(int64(20), nil)

	repoMock.On("LengthOfCollection", "db3", "collection7", mock.Anything).Return(int64(80), nil)
	repoMock.On("LengthOfCollection", "db3", "collection8", mock.Anything).Return(int64(50), nil)
	repoMock.On("LengthOfCollection", "db3", "collection9", mock.Anything).Return(int64(70), nil)

	expectResult := []model.CollectionStats{
		{CollectionName: "collection3", DocumentCount: 100},
		{CollectionName: "collection4", DocumentCount: 90},
		{CollectionName: "collection7", DocumentCount: 80},
		{CollectionName: "collection9", DocumentCount: 70},
		{CollectionName: "collection2", DocumentCount: 60},
		{CollectionName: "collection8", DocumentCount: 50},
		{CollectionName: "collection5", DocumentCount: 40},
		{CollectionName: "collection6", DocumentCount: 20},
		{CollectionName: "collection1", DocumentCount: 10},
	}

	collectionMiddleware := middleware.NewCollectionMiddleware(repoMock)

	result, err := collectionMiddleware.GetNumberOfFilesFromAllDB(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result, expectResult)
}

func TestGetListOfCollectionsStats_MaxToMin(t *testing.T) {
	repoMock := new(MockCollectionsRepository)

	expectedCollections := []string{"collection1", "collection2", "collection3"}
	repoMock.On("ListCollectionNames", "databaseName", mock.Anything).Return(expectedCollections, nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection1", mock.Anything).Return(int64(10), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection2", mock.Anything).Return(int64(20), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection3", mock.Anything).Return(int64(30), nil)

	expectResult := []model.CollectionStats{
		{
			CollectionName: "collection3",
			DocumentCount:  30,
		},
		{
			CollectionName: "collection2",
			DocumentCount:  20,
		},
		{
			CollectionName: "collection1",
			DocumentCount:  10,
		},
	}

	collectionMiddleware := middleware.NewCollectionMiddleware(repoMock)

	result, err := collectionMiddleware.GetListOfCollectionsStats("databaseName", context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result, expectResult)
}

func TestGetListOfCollectionsStats_RandomOrder(t *testing.T) {
	repoMock := new(MockCollectionsRepository)

	expectedCollections := []string{"collection1", "collection2", "collection3", "collection4", "collection5", "collection6", "collection7", "collection8", "collection9", "collection10"}
	repoMock.On("ListCollectionNames", "databaseName", mock.Anything).Return(expectedCollections, nil)

	repoMock.On("LengthOfCollection", "databaseName", "collection1", mock.Anything).Return(int64(50), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection2", mock.Anything).Return(int64(30), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection3", mock.Anything).Return(int64(20), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection4", mock.Anything).Return(int64(40), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection5", mock.Anything).Return(int64(10), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection6", mock.Anything).Return(int64(60), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection7", mock.Anything).Return(int64(100), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection8", mock.Anything).Return(int64(80), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection9", mock.Anything).Return(int64(90), nil)
	repoMock.On("LengthOfCollection", "databaseName", "collection10", mock.Anything).Return(int64(70), nil)

	expectResult := []model.CollectionStats{
		{CollectionName: "collection7", DocumentCount: 100},
		{CollectionName: "collection9", DocumentCount: 90},
		{CollectionName: "collection8", DocumentCount: 80},
		{CollectionName: "collection10", DocumentCount: 70},
		{CollectionName: "collection6", DocumentCount: 60},
		{CollectionName: "collection1", DocumentCount: 50},
		{CollectionName: "collection4", DocumentCount: 40},
		{CollectionName: "collection2", DocumentCount: 30},
		{CollectionName: "collection3", DocumentCount: 20},
		{CollectionName: "collection5", DocumentCount: 10},
	}

	collectionMiddleware := middleware.NewCollectionMiddleware(repoMock)

	result, err := collectionMiddleware.GetListOfCollectionsStats("databaseName", context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result, expectResult)
}
