package middleware

import (
	"MongoDBCounterService/internal/model"
	"log"
	"sort"
)

func reverseSortByQuantity(stats []model.CollectionStats) {
	lessFunc := func(i, j int) bool {
		return stats[i].DocumentCount > stats[j].DocumentCount
	}

	sort.Slice(stats, lessFunc)
}

func mergeReverseSort(allCollections [][]model.CollectionStats) []model.CollectionStats {
	if len(allCollections) == 0 {
		log.Printf("MongoDB is empty")
		return []model.CollectionStats{}
	}

	//todo: it can by async
	for len(allCollections) > 1 {
		var merged [][]model.CollectionStats
		for i := 0; i < len(allCollections); i += 2 {
			if i+1 < len(allCollections) {
				merged = append(merged, merge(allCollections[i], allCollections[i+1]))
			} else {
				merged = append(merged, allCollections[i])
			}
		}
		allCollections = merged
	}

	return allCollections[0]
}

func merge(list1, list2 []model.CollectionStats) []model.CollectionStats {
	mergedList := make([]model.CollectionStats, 0, len(list1)+len(list2))
	i, j := 0, 0

	for i < len(list1) && j < len(list2) {
		if list1[i].DocumentCount > list2[j].DocumentCount {
			mergedList = append(mergedList, list1[i])
			i++
		} else {
			mergedList = append(mergedList, list2[j])
			j++
		}
	}

	for ; i < len(list1); i++ {
		mergedList = append(mergedList, list1[i])
	}

	for ; j < len(list2); j++ {
		mergedList = append(mergedList, list2[j])
	}

	return mergedList
}
