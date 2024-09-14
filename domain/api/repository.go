package api

type Repository interface {
	GetAllApis() (*[]Api, error)
	GetApiByNameOrPath(name string, path string) (*Api, error)
	GetApiDetailsByNameAndPath(name string, path string) (*Api, error)
	InsertApi(apiReq *CreateApiRequest) error
	FromDomain(domainApi *Api) (any, error)
	ToDomain(repoApi any) (*Api, error)
}
