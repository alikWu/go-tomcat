package servlet

import (
	"io"
)

type ServletRequest interface {
	GetCharacterEncoding() string
	GetContentLength() int64
	GetContentType() string

	GetInputStream() io.Reader

	GetParameter(s string) string
	GetParameterMap() map[string][]string
	GetParameterNames() []string
	GetParameterValues(s string) []string

	GetProtocol() string
	GetServerPort() int64
	GetRemotePort() int64
	GetRemoteHost() string

	GetServletContext() ServletContext

	//atribute is from  sesison, cookie
	GetAttribute(name string) interface{}
	GetAttributeNames() []string
	//set attribute into session
	SetAttribute(name string, value interface{})
	RemoveAttribute(name string)
	//todo support async
}
