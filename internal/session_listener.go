package internal

type SessionListener interface {
	SessionEvent(event *SessionEvent)
}
