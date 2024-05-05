package internal

import "context"

type ValveContext interface {
	GetInfo() string
	InvokeNext(request HttpServletRequest, response HttpServletResponse) error
	GetContext() context.Context
}
