package internal

type InstanceListener interface {
	InstanceEvent(event InstanceEvent)
}
