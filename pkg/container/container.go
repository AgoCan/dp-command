package container

type Containerd struct{}

func newContainerd() *Containerd {
	return &Containerd{}
}
func (c *Containerd) Pull() {}
func (c *Containerd) Push() {}
func (c *Containerd) Tag()  {}
