package core

type FilterMap struct {
	filterName  string
	servletName string
	urlPattern  string
}

func (fm *FilterMap) GetFilterName() string {
	return fm.filterName
}

func (fm *FilterMap) SetFilterName(filterName string) {
	fm.filterName = filterName
}

func (fm *FilterMap) GetServletName() string {
	return fm.servletName
}

func (fm *FilterMap) SetServletName(servletName string) {
	fm.servletName = servletName
}

func (fm *FilterMap) GetURLPattern() string {
	return fm.urlPattern
}

func (fm *FilterMap) SetURLPatterns(urlPattern string) {
	fm.urlPattern = urlPattern
}
