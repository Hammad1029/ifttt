package rule

type Repository interface {
	GetAllRules() (*[]Rule, error)
	GetRulesByIds(ids []uint) (*[]Rule, error)
	GetRuleByName(name string) (*Rule, error)
	InsertRule(rule *CreateRuleRequest) error
}
