package servlet

type Servlet interface {
	GetServletName() string
	Service(request ServletRequest, response ServletResponse) error
	GetServletInfo() string
	GetMatchedUrlPattern() []string
}
