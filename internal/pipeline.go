package internal

import "context"

type Pipeline interface {
	GetBasic() Valve
	SetBasic(valve Valve)
	AddValve(valve Valve)
	GetValves() []Valve
	Invoke(ctx context.Context, request HttpServletRequest, response HttpServletResponse) error
	RemoveValve(valve Valve)
}
