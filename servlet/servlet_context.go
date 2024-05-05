package servlet

import "context"

type ServletContext interface {
	GetContext() context.Context
}
