package core

type ServletMap struct {
	servletName string
	urlPattern  string
}

func NewServletMap(servletName string, urlPattern string) *ServletMap {
	return &ServletMap{
		servletName: servletName,
		urlPattern:  urlPattern,
	}
}

func (sm *ServletMap) GetServletName() string {
	return sm.servletName
}

func (sm *ServletMap) SetServletName(servletName string) {
	sm.servletName = servletName
}

func (sm *ServletMap) GetURLPattern() string {
	return sm.urlPattern
}

func (sm *ServletMap) SetURLPatterns(urlPattern string) {
	sm.urlPattern = urlPattern
}
