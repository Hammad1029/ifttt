package api

type Repository interface {
	GetAllApis() (*[]ApiSerialized, error)
	GetApiByGroupAndName(group string, name string) (*ApiSerialized, bool, error)
	InsertApi(newApi *ApiSerialized) error
}
