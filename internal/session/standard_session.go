package session

import (
	"time"

	"go-tomcat/internal"
)

type standardSession struct {
	sessionID      string
	valid          bool
	attributes     map[string]interface{}
	validTimestamp int64
	creationTime   int64
	lastAccessTime int64
	listeners      []internal.SessionListener
}

func NewSession(sessionID string, expireTime int64) internal.HttpSession {
	now := time.Now().Unix()
	return &standardSession{
		sessionID:      sessionID,
		creationTime:   now,
		valid:          true,
		attributes:     make(map[string]interface{}),
		validTimestamp: now + expireTime,
		lastAccessTime: now,
	}
}

//
//func (s *standardSession) AddSessionListener(listener internal.SessionListener) {
//	s.listeners = append(s.listeners, listener)
//}
//
//func (s *standardSession) FireSessionEvent(kind string, data interface{}) {
//	event := internal.NewSessionEvent(s, kind, data)
//	for _, listener := range s.listeners {
//		listener.SessionEvent(event)
//	}
//}

func (s *standardSession) IsValid() bool {
	return s.valid
}

func (s *standardSession) GetCreationTime() int64 {
	return s.creationTime
}

func (s *standardSession) GetId() string {
	return s.sessionID
}

func (s *standardSession) GetLastAccessedTime() int64 {
	return s.lastAccessTime
}

func (s *standardSession) SetMaxInactiveInterval(arg int64) {
	//TODO implement me
	panic("implement me")
}

func (s *standardSession) GetMaxInactiveInterval() int64 {
	//TODO implement me
	panic("implement me")
}

func (s *standardSession) GetAttribute(name string) interface{} {
	if !s.IsValid() {
		return nil
	}
	return s.attributes[name]
}

func (s *standardSession) GetAttributeNames() []string {
	if !s.IsValid() {
		return []string{}
	}

	var res []string
	for name, _ := range s.attributes {
		res = append(res, name)
	}
	return res
}

func (s *standardSession) SetAttribute(name string, value interface{}) {
	if !s.IsValid() {
		return
	}
	s.attributes[name] = value
}

func (s *standardSession) RemoveAttribute(name string) {
	delete(s.attributes, name)
}

func (s *standardSession) Invalidate() {
	s.valid = false
}
