package api

type Repository interface {
	GetAllApis() (*[]Api, error)
	GetApiByNameOrPath(name string, path string) (*Api, error)
	InsertApi(apiReq *CreateApiRequest) error
}
