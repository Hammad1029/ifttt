package postgres

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/rule"

	"github.com/jackc/pgtype"
)

type PostgresRulesRepository struct {
	*PostgresBaseRepository
}

func NewPostgresRulesRepository(base *PostgresBaseRepository) *PostgresRulesRepository {
	return &PostgresRulesRepository{PostgresBaseRepository: base}
}

func (p *PostgresRulesRepository) GetAllRules() (*[]rule.Rule, error) {
	var pgRules []rules
	if err := p.client.Find(&pgRules).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetAllRules: could not query rules: %s", err)
	}

	var domainRules []rule.Rule
	for _, r := range pgRules {
		if dRule, err := p.ToDomain(&r); err != nil {
			return nil,
				fmt.Errorf("method *PostgresApiRepository.GetAllRules: could not convert to domain rule: %s", err)
		} else {
			domainRules = append(domainRules, *dRule)
		}
	}

	return &domainRules, nil
}

func (p *PostgresRulesRepository) GetRulesByIds(ids []uint) (*[]rule.Rule, error) {
	var pgRules []rules
	if err := p.client.Find(&pgRules, "ids in ?", ids).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetRulesByIds: could not query rules: %s", err)
	}

	var domainRules []rule.Rule
	for _, r := range pgRules {
		if dRule, err := p.ToDomain(&r); err != nil {
			return nil,
				fmt.Errorf("method *PostgresApiRepository.GetAllRules: could not convert to domain rule: %s", err)
		} else {
			domainRules = append(domainRules, *dRule)
		}
	}

	return &domainRules, nil
}

func (p *PostgresRulesRepository) GetRuleByName(name string) (*rule.Rule, error) {
	var pgRule rules
	if err := p.client.First(&pgRule, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetRuleByName: could not query rules: %s", err)
	}

	var domainRule rule.Rule
	if dRule, err := p.ToDomain(&pgRule); err != nil {
		return nil,
			fmt.Errorf("method *PostgresApiRepository.GetAllRules: could not convert to domain rule: %s", err)
	} else {
		domainRule = *dRule
	}

	return &domainRule, nil
}

func (p *PostgresRulesRepository) InsertRule(dRule *rule.CreateRuleRequest) error {
	pgRule, err := p.FromDomain(&rule.Rule{
		Name:        dRule.Name,
		Description: dRule.Description,
		Conditions:  dRule.Conditions,
		Then:        dRule.Then,
		Else:        dRule.Else,
	})
	if err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertRule: could not convert to pgRule: %s", err)
	}

	if err := p.client.Create(&pgRule).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertRule: could not insert rules: %s", err)
	}

	return nil
}

func (p *PostgresRulesRepository) FromDomain(domainRule *rule.Rule) (any, error) {
	pgRule := rules{
		Name:        domainRule.Name,
		Description: domainRule.Description,
	}

	if conditionsMarshalled, err := json.Marshal(domainRule.Conditions); err != nil {
		return nil, fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal conditions: %s", err)
	} else {
		pgRule.Conditions = pgtype.JSONB{Bytes: conditionsMarshalled}
	}

	if thenMarshalled, err := json.Marshal(domainRule.Then); err != nil {
		return nil, fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal then: %s", err)
	} else {
		pgRule.Then = pgtype.JSONB{Bytes: thenMarshalled}
	}

	if elseMarshalled, err := json.Marshal(domainRule.Else); err != nil {
		return nil, fmt.Errorf("method *PostgresRulesRepository.FromDomain: could not marshal else: %s", err)
	} else {
		pgRule.Else = pgtype.JSONB{Bytes: elseMarshalled}
	}

	return &pgRule, nil
}

func (p *PostgresRulesRepository) ToDomain(repoRule any) (*rule.Rule, error) {
	pgRule, ok := repoRule.(*rules)
	if !ok {
		return nil, fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in casting pgRule")
	}

	domainRule := rule.Rule{
		Name:        pgRule.Name,
		Description: pgRule.Description,
	}

	if err := json.Unmarshal(pgRule.Conditions.Bytes, &domainRule.Conditions); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling conditions: %s", err)
	}

	if err := json.Unmarshal(pgRule.Then.Bytes, &domainRule.Then); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling then: %s", err)
	}

	if err := json.Unmarshal(pgRule.Else.Bytes, &domainRule.Else); err != nil {
		return nil,
			fmt.Errorf("method *PostgresRulesRepository.ToDomain: error in unmarshalling else: %s", err)
	}

	return &domainRule, nil
}
