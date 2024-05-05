package connector

import (
	"io"

	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/cookie"
	"github.com/alikWu/go-tomcat/servlet"
)

/**
reason: 直接使用的是 HttpRequest 与 HttpResponse，这两个对象要传入 Servlet 中，但在这两个类中我们也定义了许多内部的方法，一旦被用户知晓我们的实现类，
则这些内部方法就暴露在用户面前了。这样其实是不好的。面向对象编程的思想，是将内部实现的结构和复杂性包装在一层壳里面，能不对外暴露就不要对外暴露。
作为客户程序，最好只知道最小知识集。另外，这个 Request 和 Response 类是要传给外部的 Servlet 程序的，跳出了 Tomcat 本身，
如果这个写 Servlet 的程序员他知道传的这个类里面有一些额外的方法，原理上他可以进行强制转换之后调用，这样也不是很安全。

aim: Facade 类是一个新类，外部使用者没法根据它来强制转换获得内部的结构和方法，这样将实际实现的几个类保护起来了，提高了安全性。
*/
type HttpRequestFacade struct {
	request internal.HttpServletRequest
}

func NewHttpRequestFacade(request internal.HttpServletRequest) internal.HttpServletRequest {
	return &HttpRequestFacade{
		request: request,
	}
}

func (h HttpRequestFacade) GetProtocol() string {
	return h.request.GetProtocol()
}

func (h HttpRequestFacade) GetServerPort() int64 {
	return h.request.GetServerPort()
}

func (h HttpRequestFacade) GetRemotePort() int64 {
	return h.request.GetRemotePort()
}

func (h HttpRequestFacade) GetRemoteHost() string {
	return h.request.GetRemoteHost()
}

func (h HttpRequestFacade) GetPathInfo() string {
	return h.request.GetPathInfo()
}

func (h HttpRequestFacade) GetSession(create bool) internal.HttpSession {
	return h.request.GetSession(create)
}

func (h HttpRequestFacade) GetCharacterEncoding() string {
	return h.request.GetCharacterEncoding()
}

func (h HttpRequestFacade) GetContentLength() int64 {
	return h.request.GetContentLength()
}

func (h HttpRequestFacade) GetContentType() string {
	return h.request.GetContentType()
}

func (h HttpRequestFacade) GetInputStream() io.Reader {
	return h.request.GetInputStream()
}

func (h HttpRequestFacade) GetParameter(s string) string {
	return h.request.GetParameter(s)
}

func (h HttpRequestFacade) GetParameterMap() map[string][]string {
	return h.request.GetParameterMap()
}

func (h HttpRequestFacade) GetParameterNames() []string {
	return h.request.GetParameterNames()
}

func (h HttpRequestFacade) GetParameterValues(s string) []string {
	return h.request.GetParameterValues(s)
}

func (h HttpRequestFacade) GetHeader(arg string) string {
	return h.request.GetHeader(arg)
}

func (h HttpRequestFacade) GetHeaderNames() []string {
	return h.request.GetHeaderNames()
}

func (h HttpRequestFacade) GetServletContext() servlet.ServletContext {
	return h.request.GetServletContext()
}

func (h HttpRequestFacade) GetCookies() []*cookie.Cookie {
	return h.request.GetCookies()
}

func (h HttpRequestFacade) GetQueryString() string {
	return h.request.GetQueryString()
}

func (h HttpRequestFacade) GetMethod() string {
	return h.request.GetMethod()
}

func (h HttpRequestFacade) GetAttribute(name string) interface{} {
	return h.request.GetAttribute(name)
}

func (h HttpRequestFacade) GetAttributeNames() []string {
	return h.request.GetAttributeNames()
}

func (h HttpRequestFacade) SetAttribute(name string, value interface{}) {
	h.request.SetAttribute(name, value)
}

func (h HttpRequestFacade) RemoveAttribute(name string) {
	h.request.RemoveAttribute(name)
}
