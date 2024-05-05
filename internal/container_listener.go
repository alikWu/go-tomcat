package internal

type ContainerListener interface {
	ContainerEvent(event *ContainerEvent)
}
