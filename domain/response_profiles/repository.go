package responseprofiles

type Repository interface {
	AddProfile(p *Profile, parent uint) error
	GetProfilesByInternalAndCode(internal bool, code string) (*[]Profile, error)
	GetAllInternalProfiles() (*[]Profile, error)
}
