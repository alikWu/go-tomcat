package internal

type Connector interface {
	GetContainer() Context
	SetContainer(container Context)
	GetProtocol() string
	SetProtocol(protocol string)
	Initialize()
	ListenConnect() error
}
