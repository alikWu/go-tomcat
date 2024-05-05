package internal

const (
	SESSION_CREATED_EVENT   = "createSession"
	SESSION_DESTROYED_EVENT = "destroySession"
)

type HttpSession interface {
	GetCreationTime() int64
	GetId() string
	GetLastAccessedTime() int64
	GetMaxInactiveInterval() int64
	SetMaxInactiveInterval(interval int64)

	GetAttribute(name string) interface{}
	GetAttributeNames() []string
	SetAttribute(name string, value interface{})
	RemoveAttribute(name string)

	IsValid() bool
	Invalidate()
}
