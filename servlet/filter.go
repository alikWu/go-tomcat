package servlet

type Filter interface {
	GetFilterName() string
	DoFilter(request ServletRequest, response ServletResponse, chain FilterChain) error
	GetFilterMatch() []*FilterMatch
}

//匹配servletName or 匹配urlPattern
type FilterMatch struct {
	servletName string
	urlPattern  string
}

func NewFilterMatchServlet(servletName string) *FilterMatch {
	return &FilterMatch{
		servletName: servletName,
	}
}

func NewFilterMatchUrl(urlPattern string) *FilterMatch {
	return &FilterMatch{
		urlPattern: urlPattern,
	}
}

func (fm *FilterMatch) GetServletName() string {
	return fm.servletName
}

func (fm *FilterMatch) GetUrlPattern() string {
	return fm.urlPattern
}
