package internal

type Valve interface {
	GetInfo() string
	GetContainer() Container
	SetContainer(container Container)
	Invoke(request HttpServletRequest, response HttpServletResponse, ctx ValveContext) error
}
