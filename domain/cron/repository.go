package cron

type Repository interface {
	GetAllCrons() (*[]Cron, error)
	GetCronByName(name string) (*Cron, error)
	InsertCron(req *Cron, apiID uint) error
}
