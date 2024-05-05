package internal

import "github.com/alikWu/go-tomcat/servlet"

type Wrapper interface {
	Container
	GetServlet() servlet.Servlet
	Load()
}
