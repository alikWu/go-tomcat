package internal

import (
	"context"
)

type Container interface {
	GetName() string
	SetName(name string)
	GetParent() Container
	SetParent(parent Container)
	AddChild(child Container)
	FindChild(name string) Container
	Invoke(ctx context.Context, request HttpServletRequest, response HttpServletResponse) error
}
