package session

import (
	"go-tomcat/internal"
)

type standardSessionFacade struct {
	hs internal.HttpSession
}

func NewSessionFacade(hs internal.HttpSession) internal.HttpSession {
	return &standardSessionFacade{hs: hs}
}

func (ssf standardSessionFacade) IsValid() bool {
	return ssf.hs.IsValid()
}

func (ssf standardSessionFacade) GetCreationTime() int64 {
	return ssf.hs.GetCreationTime()
}

func (ssf standardSessionFacade) GetId() string {
	return ssf.hs.GetId()
}

func (ssf standardSessionFacade) GetLastAccessedTime() int64 {
	return ssf.hs.GetLastAccessedTime()
}

func (ssf standardSessionFacade) SetMaxInactiveInterval(arg int64) {
	ssf.hs.SetMaxInactiveInterval(arg)
}

func (ssf standardSessionFacade) GetMaxInactiveInterval() int64 {
	return ssf.hs.GetMaxInactiveInterval()
}

func (ssf standardSessionFacade) GetAttribute(name string) interface{} {
	return ssf.hs.GetAttribute(name)
}

func (ssf standardSessionFacade) GetAttributeNames() []string {
	return ssf.hs.GetAttributeNames()
}

func (ssf standardSessionFacade) SetAttribute(name string, value interface{}) {
	ssf.hs.SetAttribute(name, value)
}

func (ssf standardSessionFacade) RemoveAttribute(name string) {
	ssf.hs.RemoveAttribute(name)
}

func (ssf standardSessionFacade) Invalidate() {
	ssf.hs.Invalidate()
}
