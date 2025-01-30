package configuration

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func ValidateMapWithInternalTags(val any, configRepo *InternalTagRepository) error {
	reflected := reflect.Indirect(reflect.ValueOf(val))
	indirectValue := reflected.Interface()
	if reflected.Kind() == reflect.Map {
		var internalTag InternalTagInMap
		if err := mapstructure.Decode(indirectValue, &internalTag); err != nil {
			return err
		} else if internalTag.InternalTag != "" {
			if pTag, err := (*configRepo).GetByIDOrName(0, internalTag.InternalTag); err != nil {
				return err
			} else if pTag == nil {
				return fmt.Errorf("tag %s not found", internalTag.InternalTag)
			}
		} else {
			mapCloned := map[string]any{}
			if err := mapstructure.Decode(indirectValue, &mapCloned); err != nil {
				return err
			}
			for _, v := range mapCloned {
				if err := ValidateMapWithInternalTags(v, configRepo); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
