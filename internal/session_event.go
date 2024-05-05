package internal

type SessionEvent struct {
	data    interface{}
	session HttpSession
	kind    string
}

func NewSessionEvent(session HttpSession, kind string, data interface{}) *SessionEvent {
	return &SessionEvent{
		data:    data,
		session: session,
		kind:    kind,
	}
}

func (se *SessionEvent) GetData() interface{} {
	return se.data
}

func (se *SessionEvent) GetSession() HttpSession {
	return se.session
}

func (se *SessionEvent) GetType() string {
	return se.kind
}
