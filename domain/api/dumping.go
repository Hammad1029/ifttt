package api

import "ifttt/manager/domain/resolvable"

type Dumping struct {
	Table    string                           `json:"table" mapstructure:"table"`
	Mappings map[string]resolvable.Resolvable `json:"mappings" mapstructure:"mappings"`
}
