package postgres

import (
	"fmt"
	"ifttt/manager/domain/cron"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"gorm.io/gorm"
)

type PostgresCronRepository struct {
	*PostgresBaseRepository
}

func NewPostgresCronRepository(base *PostgresBaseRepository) *PostgresCronRepository {
	return &PostgresCronRepository{PostgresBaseRepository: base}
}

func (p *PostgresCronRepository) GetAllCrons() (*[]cron.Cron, error) {
	var pgCron []crons
	if err := p.client.
		Preload("TriggerFlowRef").Preload("TriggerFlowRef.Rules").
		Find(&pgCron).Error; err != nil {
		return nil, err
	}

	var domainCrons []cron.Cron
	for _, c := range pgCron {
		if dCron, err := c.toDomain(); err != nil {
			return nil, err
		} else {
			domainCrons = append(domainCrons, *dCron)
		}
	}

	return &domainCrons, nil
}

func (p *PostgresCronRepository) GetCronByName(name string) (*cron.Cron, error) {
	var pgCron crons
	if err := p.client.
		Preload("TriggerFlowRef").Preload("TriggerFlowRef.Rules").
		First(&pgCron, "name ilike ?", fmt.Sprintf("%%%s%%", name)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var domainCron cron.Cron
	if dCron, err := pgCron.toDomain(); err != nil {
		return nil, err
	} else {
		domainCron = *dCron
	}

	return &domainCron, nil
}

func (p *PostgresCronRepository) InsertCron(req *cron.CreateCronRequest, attachTriggers *[]triggerflow.TriggerFlow) error {
	var pgCron crons
	err := pgCron.fromDomain(req, attachTriggers)
	if err != nil {
		return err
	}

	if err := p.client.Create(&pgCron).Error; err != nil {
		return err
	}
	return nil
}
