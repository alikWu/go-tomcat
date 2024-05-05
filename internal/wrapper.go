package internal

import "go-tomcat/servlet"

type Wrapper interface {
	Container
	GetServlet() servlet.Servlet
	Load()
}
