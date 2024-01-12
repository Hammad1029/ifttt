package schemas

type AddApi struct {
	ClientId string `json:"clientId" binding:"required"`
	ApiName  string `json:"apiName" binding:"required"`
	PathName string `json:"pathName" binding:"required"`
}

type AddMappingToApi struct {
	ClientId string            `json:"clientId" binding:"required"`
	ApiName  string            `json:"apiName" binding:"required"`
	Mappings map[string]string `json:"mappings" binding:"required"`
}
