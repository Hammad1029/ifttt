package api

type Repository interface {
	GetAllApis() (*[]Api, error)
	GetApiByNameAndPath(name string, path string) (*Api, error)
	InsertApi(apiReq *CreateApiRequest) error
}
