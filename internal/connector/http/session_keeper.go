package http

import (
	"sync"

	"github.com/google/uuid"
	"go-tomcat/internal"
	"go-tomcat/internal/session"
)

type SessionKeeper struct {
	mutex      sync.Mutex
	sessionMap map[string]internal.HttpSession
}

func NewSessionKeeper() *SessionKeeper {
	return &SessionKeeper{
		sessionMap: make(map[string]internal.HttpSession, 256),
	}
}

func (sk *SessionKeeper) CreateSession(expireTime int64) internal.HttpSession {
	sessionID := uuid.New().String()
	ses := session.NewSession(sessionID, expireTime)

	sk.mutex.Lock()
	defer func() {
		sk.mutex.Unlock()
	}()

	sk.sessionMap[sessionID] = ses
	return ses
}

func (sk *SessionKeeper) GetSession(sessionId string) internal.HttpSession {
	httpSession := sk.sessionMap[sessionId]
	if httpSession == nil {
		return nil
	}

	sk.mutex.Lock()
	defer func() {
		sk.mutex.Unlock()
	}()
	if !httpSession.IsValid() {
		delete(sk.sessionMap, sessionId)
		return nil
	}
	return httpSession
}

func (sk *SessionKeeper) DeprecateInvalidSessions() {
	var invalidSessionID []string
	for sessionID, httpSession := range sk.sessionMap {
		if !httpSession.IsValid() {
			invalidSessionID = append(invalidSessionID, sessionID)
		}
	}

	sk.mutex.Lock()
	defer func() {
		sk.mutex.Unlock()
	}()
	for _, sessionID := range invalidSessionID {
		delete(sk.sessionMap, sessionID)
	}
}
