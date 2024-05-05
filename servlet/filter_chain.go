package servlet

type FilterChain interface {
	DoFilter(request ServletRequest, response ServletResponse) error
}
