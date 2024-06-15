package schemas

type AddIndexRequest struct {
	TableName      string   `json:"tableName"`
	IndexedColumns []string `json:"indexedColumns"`
	Local          bool     `json:"local"`
}

type FindIndexRequest struct {
	IndexedColumns []string `json:"indexedColumns"`
	TableName      string   `json:"tableName"`
}

type DropIndexRequest struct {
	IndexName string `json:"indexName"`
	TableName string `json:"tableName"`
}
