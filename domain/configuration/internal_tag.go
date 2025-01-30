package configuration

type InternalTagGroup struct {
	ID   uint           `json:"id" mapstructure:"id"`
	Name string         `json:"name" mapstructure:"name"`
	Tags []*InternalTag `json:"tags" mapstructure:"tags"`
}

type InternalTag struct {
	ID       uint               `json:"id" mapstructure:"id"`
	Name     string             `json:"name" mapstructure:"name"`
	Groups   []InternalTagGroup `json:"groups" mapstructure:"groups"`
	All      bool               `json:"all" mapstructure:"all"`
	Reserved bool               `json:"reserved" mapstructure:"reserved"`
}

type InternalTagInMap struct {
	InternalTag string `json:"internalTag" mapstructure:"internalTag"`
}

type InternalTagRepository interface {
	AddGroup(g *InternalTagGroup) error
	GetAllGroups() (*[]InternalTagGroup, error)
	GetGroupByName(name string) (*InternalTagGroup, error)
	Add(pTag *InternalTag) error
	GetByIDOrName(id uint, name string) (*InternalTag, error)
	GetAll() (*[]InternalTag, error)
}
