package internal

type ListenerFactory interface {
	CreateContainerListener(listenerName string) ContainerListener
}
