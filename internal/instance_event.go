package internal

import "github.com/alikWu/go-tomcat/servlet"

const (
	BEFORE_INIT_EVENT     = "beforeInit"
	AFTER_INIT_EVENT      = "afterInit"
	BEFORE_SERVICE_EVENT  = "beforeService"
	AFTER_SERVICE_EVENT   = "afterService"
	BEFORE_DESTROY_EVENT  = "beforeDestroy"
	AFTER_DESTROY_EVENT   = "afterDestroy"
	BEFORE_DISPATCH_EVENT = "beforeDispatch"
	AFTER_DISPATCH_EVENT  = "afterDispatch"
	BEFORE_FILTER_EVENT   = "beforeFilter"
	AFTER_FILTER_EVENT    = "afterFilter"
)

type InstanceEvent struct {
	filter   servlet.Filter
	request  servlet.ServletRequest
	response servlet.ServletResponse
	servlet  servlet.Servlet
	kind     string
	wrapper  Wrapper
	err      error
}

func (ie *InstanceEvent) GetFilter() servlet.Filter {
	return ie.filter
}

func (ie *InstanceEvent) GetRequest() servlet.ServletRequest {
	return ie.request
}

func (ie *InstanceEvent) GetResponse() servlet.ServletResponse {
	return ie.response
}

func (ie *InstanceEvent) GetServlet() servlet.Servlet {
	return ie.servlet
}

func (ie *InstanceEvent) GetType() string {
	return ie.kind
}

func (ie *InstanceEvent) GetWrapper() Wrapper {
	return ie.wrapper
}

func (ie *InstanceEvent) GetError() error {
	return ie.err
}
