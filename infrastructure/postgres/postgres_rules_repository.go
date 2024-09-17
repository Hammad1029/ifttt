package postgres

import (
	"fmt"
	"ifttt/manager/domain/rule"

	"gorm.io/gorm"
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
		if dRule, err := r.toDomain(); err != nil {
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
	if err := p.client.Find(&pgRules, "id in ?", ids).Error; err != nil {
		return nil, fmt.Errorf("method *PostgresApiRepository.GetRulesByIds: could not query rules: %s", err)
	}

	var domainRules []rule.Rule
	for _, r := range pgRules {
		if dRule, err := r.toDomain(); err != nil {
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("method *PostgresApiRepository.GetRuleByName: could not query rules: %s", err)
	}

	var domainRule rule.Rule
	if dRule, err := pgRule.toDomain(); err != nil {
		return nil,
			fmt.Errorf("method *PostgresApiRepository.GetAllRules: could not convert to domain rule: %s", err)
	} else {
		domainRule = *dRule
	}

	return &domainRule, nil
}

func (p *PostgresRulesRepository) InsertRule(dRule *rule.CreateRuleRequest) error {
	var pgRule rules
	if err := pgRule.fromDomain(dRule); err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertRule: could not convert to pgRule: %s", err)
	}

	if err := p.client.Create(&pgRule).Error; err != nil {
		return fmt.Errorf("method *PostgresApiRepository.InsertRule: could not insert rules: %s", err)
	}

	return nil
}
