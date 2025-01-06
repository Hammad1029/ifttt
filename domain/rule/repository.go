package rule

type Repository interface {
	GetAllRules() (*[]Rule, error)
	GetRulesByIds(ids []uint) (*[]Rule, error)
	GetRulesByNames(names []string) (*[]Rule, error)
	GetRuleByName(name string) (*Rule, error)
	GetRulesLikeName(name string) (*[]Rule, error)
	InsertRule(rule *CreateRuleRequest) error
}
