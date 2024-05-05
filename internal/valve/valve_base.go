package valve

import "github.com/alikWu/go-tomcat/internal"

type ValveBase struct {
	container internal.Container
	info      string
}

func (vb *ValveBase) GetContainer() internal.Container {
	return vb.container
}

func (vb *ValveBase) SetContainer(container internal.Container) {
	vb.container = container
}

func (vb *ValveBase) GetInfo() string {
	return "ValveBase 0.1"
}
