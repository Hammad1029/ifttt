package eventprofiles

type Repository interface {
	AddProfile(p *Profile, parent uint) error
	GetProfilesByInternalAndTrigger(internal bool, trigger string) (*[]Profile, error)
	GetAllInternalProfiles() (*[]Profile, error)
}
