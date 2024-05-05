package internal

import (
	"github.com/alikWu/go-tomcat/internal/cookie"
	"github.com/alikWu/go-tomcat/servlet"
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
