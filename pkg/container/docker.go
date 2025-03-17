package container

type Docker struct{}

func newDocker() *Docker {
	return &Docker{}
}
func (d *Docker) Pull() {}
func (d *Docker) Push() {}
func (d *Docker) Tag()  {}
