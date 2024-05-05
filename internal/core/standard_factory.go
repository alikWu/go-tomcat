package core

import (
	"go-tomcat/internal"
	"go-tomcat/servlet"
	"go-tomcat/webapps/test"
)

type standardFactory struct {
	servlets           map[string]servlet.Servlet
	filters            map[string]servlet.Filter
	containerListeners map[string]internal.ContainerListener
}

func NewStandardFactory() *standardFactory {
	return &standardFactory{
		servlets:           make(map[string]servlet.Servlet),
		filters:            make(map[string]servlet.Filter),
		containerListeners: make(map[string]internal.ContainerListener),
	}
}

func (sf *standardFactory) CreateContainerListener(listenerName string) internal.ContainerListener {
	return &test.TestListener{}
}
