package api

import (
	"fmt"
	"ifttt/manager/domain/configuration"
)

func ValidateResponseDefinition(
	def *ResponseDefinition, profileRepo *configuration.ResponseProfileRepository, tagRepo *configuration.InternalTagRepository,
) error {
	if def.UseProfile == "" {
		if err := configuration.ValidateMapWithInternalTags(def.Definition, tagRepo); err != nil {
			return err
		}
	} else {
		if existing, err := (*profileRepo).GetProfilesByName(def.UseProfile); err != nil {
			return err
		} else if existing == nil || len(*existing) == 0 {
			return fmt.Errorf("response profile %s not found", def.UseProfile)
		}
	}
	return nil
}
