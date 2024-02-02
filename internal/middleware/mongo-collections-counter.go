package middleware

import (
	"MongoDBCounterService/internal/model"
	"MongoDBCounterService/internal/repository"
	"context"
	"log"
	"sync"
)

type Collection interface {
	GetNumberOfFilesFromAllDB(ctx context.Context) ([]model.CollectionStats, error)
	GetListOfCollectionsStats(database string, ctx context.Context) ([]model.CollectionStats, error)
}

type CollectionMiddleware struct {
	rep repository.CollectionsRepository
}

func NewCollectionMiddleware(collectionsRepository repository.CollectionsRepository) Collection {
	return &CollectionMiddleware{rep: collectionsRepository}
}

type CollectionTask struct {
	database   string
	collection string
	ctx        context.Context
	stats      *[]model.CollectionStats
	mu         *sync.Mutex
}

func (m *CollectionMiddleware) GetNumberOfFilesFromAllDB(ctx context.Context) ([]model.CollectionStats, error) {
	databases, err := m.rep.ListDatabaseNames(ctx)
	if err != nil {
		log.Printf("Can't get list of database names, %v", err)
		return nil, err
	}

	var allCollections [][]model.CollectionStats
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, database := range databases {
		wg.Add(1)
		go func(database string) {
			defer wg.Done()
			collectionsList, err := m.GetListOfCollectionsStats(database, ctx)
			if err != nil {
				log.Printf("Can't get list of CollectionsStats for %s, err: %v", database, err)
				return
			}

			mu.Lock()
			allCollections = append(allCollections, collectionsList)
			mu.Unlock()

		}(database)
	}

	wg.Wait()
	finalListOfCollections := mergeSort(allCollections)
	return finalListOfCollections, nil
}

func (m *CollectionMiddleware) GetListOfCollectionsStats(database string, ctx context.Context) ([]model.CollectionStats, error) {
	collections, err := m.rep.ListCollectionNames(database, ctx)
	if err != nil {
		log.Printf("Can't get list of collections names, %v", err)
		return nil, err
	}

	var stats []model.CollectionStats
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, collection := range collections {
		wg.Add(1)
		task := CollectionTask{
			database:   database,
			collection: collection,
			ctx:        ctx,
			stats:      &stats,
			mu:         &mu,
		}
		go processCollection(task, &wg, m.rep)
	}

	wg.Wait()
	sortByQuantity(stats)
	return stats, nil
}

func processCollection(task CollectionTask, wg *sync.WaitGroup, db repository.CollectionsRepository) {
	defer wg.Done()

	count, err := db.LengthOfCollection(task.database, task.collection, task.ctx)
	if err != nil {
		log.Printf("Can't count files in %s, err: %v", task.collection, err)
		count = 0
	}

	task.mu.Lock()
	*task.stats = append(*task.stats, model.CollectionStats{
		CollectionName: task.collection,
		DocumentCount:  count,
	})
	task.mu.Unlock()
}
