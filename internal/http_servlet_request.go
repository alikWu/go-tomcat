package internal

import (
	"go-tomcat/internal/cookie"
	"go-tomcat/servlet"
)

type HttpServletRequest interface {
	GetMethod() string
	GetPathInfo() string
	GetQueryString() string

	GetHeader(arg string) string
	GetHeaderNames() []string

	GetCookies() []*cookie.Cookie
	GetSession(create bool) HttpSession

	servlet.ServletRequest
}
