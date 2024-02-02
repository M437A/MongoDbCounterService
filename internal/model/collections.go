package model

type CollectionStats struct {
	CollectionName string
	DocumentCount  int64
}

type DataBaseName struct {
	DataBaseName string `json:"database"`
}
