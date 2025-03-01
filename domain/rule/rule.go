package rule

import (
	"ifttt/manager/domain/resolvable"
)

type Rule struct {
	ID          uint                    `json:"id" mapstructure:"id"`
	Name        string                  `json:"name" mapstructure:"name"`
	Description string                  `json:"description" mapstructure:"description"`
	Pre         []resolvable.Resolvable `json:"pre" mapstructure:"pre"`
	Switch      RuleSwitch              `json:"switch" mapstructure:"switch"`
	Finally     []resolvable.Resolvable `json:"finally" mapstructure:"finally"`
}

type RuleSwitch struct {
	Cases   []RuleSwitchCase `json:"cases" mapstructure:"cases"`
	Default RuleSwitchCase   `json:"default" mapstructure:"default"`
}

type RuleSwitchCase struct {
	Condition resolvable.Condition    `json:"condition" mapstructure:"condition"`
	Do        []resolvable.Resolvable `json:"do" mapstructure:"do"`
	Return    uint                    `json:"return" mapstructure:"return"`
}
