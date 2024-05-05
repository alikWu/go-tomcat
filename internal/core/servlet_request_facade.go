package core

import (
	"io"

	"github.com/alikWu/go-tomcat/servlet"
)

type servletRequestFacade struct {
	request servlet.ServletRequest
}

func NewServletRequestFacade(request servlet.ServletRequest) servlet.ServletRequest {
	return &servletRequestFacade{
		request: request,
	}
}

func (srf servletRequestFacade) GetCharacterEncoding() string {
	return srf.request.GetCharacterEncoding()
}

func (srf servletRequestFacade) GetContentLength() int64 {
	return srf.request.GetContentLength()
}

func (srf servletRequestFacade) GetContentType() string {
	return srf.request.GetContentType()
}

func (srf servletRequestFacade) GetInputStream() io.Reader {
	return srf.request.GetInputStream()
}

func (srf servletRequestFacade) GetParameter(s string) string {
	return srf.request.GetParameter(s)
}

func (srf servletRequestFacade) GetParameterMap() map[string][]string {
	return srf.request.GetParameterMap()
}

func (srf servletRequestFacade) GetParameterNames() []string {
	return srf.request.GetParameterNames()
}

func (srf servletRequestFacade) GetParameterValues(s string) []string {
	return srf.request.GetParameterValues(s)
}

func (srf servletRequestFacade) GetProtocol() string {
	return srf.request.GetProtocol()
}

func (srf servletRequestFacade) GetServerPort() int64 {
	return srf.request.GetServerPort()
}

func (srf servletRequestFacade) GetRemotePort() int64 {
	return srf.request.GetRemotePort()
}

func (srf servletRequestFacade) GetRemoteHost() string {
	return srf.request.GetRemoteHost()
}

func (srf servletRequestFacade) GetServletContext() servlet.ServletContext {
	return srf.request.GetServletContext()
}

func (srf servletRequestFacade) GetAttribute(name string) interface{} {
	return srf.request.GetAttribute(name)
}

func (srf servletRequestFacade) GetAttributeNames() []string {
	return srf.request.GetAttributeNames()
}

func (srf servletRequestFacade) SetAttribute(name string, value interface{}) {
	srf.request.SetAttribute(name, value)
}

func (srf servletRequestFacade) RemoveAttribute(name string) {
	srf.request.RemoveAttribute(name)
}
