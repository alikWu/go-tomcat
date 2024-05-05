package internal

type ContainerEvent struct {
	container Container
	data      interface{}
	//type
	kind string
}

func NewContainerEvent(container Container, kind string, data interface{}) *ContainerEvent {
	return &ContainerEvent{
		container: container,
		data:      data,
		kind:      kind,
	}
}

func (ce *ContainerEvent) GetData() interface{} {
	return ce.data
}

func (ce *ContainerEvent) GetContainer() Container {
	return ce.container
}

func (ce *ContainerEvent) GetType() string {
	return ce.kind
}

func (ce *ContainerEvent) ToString() string {
	return "ContainerEvent[]"
}
