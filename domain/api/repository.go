package api

type Repository interface {
	ToLocal(input *Api, output any) error
	ToGlobal(input any) (*Api, error)
	GetAllApis() (*[]Api, error)
	GetApisByGroupAndName(group string, name string) (*[]Api, error)
	InsertApi(newApi *Api) error
}
