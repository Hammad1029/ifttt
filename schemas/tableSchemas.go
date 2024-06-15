package schemas

import "generic/models"

type AddTableRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Columns     []models.TableColumn `json:"columns"`
}

type GetTablesRequest struct {
	InternalName string `json:"internalName"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type ColumnChange struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Add      bool   `json:"add"`
	Drop     bool   `json:"drop"`
	Alter    bool   `json:"alter"`
	Rename   string `json:"rename"`
}

type UpdateTableRequest struct {
	InternalName string         `json:"internalName"`
	Description  string         `json:"description"`
	Columns      []ColumnChange `json:"columns"`
}
