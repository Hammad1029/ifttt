package requestvalidator

import "ifttt/manager/domain/configuration"

func GenerateAll(req *map[string]RequestParameter, configRepo *configuration.InternalTagRepository) error {
	tagsInUse := make(map[uint]bool)
	for key, val := range *req {
		if regex, err := val.Generate(configRepo, &tagsInUse); err != nil {
			return err
		} else {
			val.Regex = regex
			(*req)[key] = val
		}
	}
	return nil
}
