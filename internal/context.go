package internal

import "go-tomcat/servlet"

//Tomcat 内的 Context 组件跟 Servlet 规范中的 ServletContext 接口的区别：
//Servlet规范中ServletContext表示web应用的上下文环境
//而web应用对应tomcat的概念是Context，所以从设计上，ServletContext自然会成为tomcat的Context具体实现的一个成员变量。
type Context interface {
	Container
	SetConnector(connector Connector)
	GetConnector() Connector

	GetServletContext() servlet.ServletContext
	GetSessionTimeout() int64
	SetSessionTimeout(timeout int64)

	GetWrapper(url string) Container
	Reload()
}
