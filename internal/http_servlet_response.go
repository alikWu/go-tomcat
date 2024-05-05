package internal

import (
	"go-tomcat/internal/cookie"
	"go-tomcat/servlet"
)

type StatusCode int32

const (
	SC_CONTINUE                        StatusCode = 100
	SC_SWITCHING_PROTOCOLS             StatusCode = 101
	SC_OK                              StatusCode = 200
	SC_CREATED                         StatusCode = 201
	SC_ACCEPTED                        StatusCode = 202
	SC_NON_AUTHORITATIVE_INFORMATION   StatusCode = 203
	SC_NO_CONTENT                      StatusCode = 204
	SC_RESET_CONTENT                   StatusCode = 205
	SC_PARTIAL_CONTENT                 StatusCode = 206
	SC_MULTIPLE_CHOICES                StatusCode = 300
	SC_MOVED_PERMANENTLY               StatusCode = 301
	SC_MOVED_TEMPORARILY               StatusCode = 302
	SC_FOUND                           StatusCode = 302
	SC_SEE_OTHER                       StatusCode = 303
	SC_NOT_MODIFIED                    StatusCode = 304
	SC_USE_PROXY                       StatusCode = 305
	SC_TEMPORARY_REDIRECT              StatusCode = 307
	SC_BAD_REQUEST                     StatusCode = 400
	SC_UNAUTHORIZED                    StatusCode = 401
	SC_PAYMENT_REQUIRED                StatusCode = 402
	SC_FORBIDDEN                       StatusCode = 403
	SC_NOT_FOUND                       StatusCode = 404
	SC_METHOD_NOT_ALLOWED              StatusCode = 405
	SC_NOT_ACCEPTABLE                  StatusCode = 406
	SC_PROXY_AUTHENTICATION_REQUIRED   StatusCode = 407
	SC_REQUEST_TIMEOUT                 StatusCode = 408
	SC_CONFLICT                        StatusCode = 409
	SC_GONE                            StatusCode = 410
	SC_LENGTH_REQUIRED                 StatusCode = 411
	SC_PRECONDITION_FAILED             StatusCode = 412
	SC_REQUEST_ENTITY_TOO_LARGE        StatusCode = 413
	SC_REQUEST_URI_TOO_LONG            StatusCode = 414
	SC_UNSUPPORTED_MEDIA_TYPE          StatusCode = 415
	SC_REQUESTED_RANGE_NOT_SATISFIABLE StatusCode = 416
	SC_EXPECTATION_FAILED              StatusCode = 417
	SC_INTERNAL_SERVER_ERROR           StatusCode = 500
	SC_NOT_IMPLEMENTED                 StatusCode = 501
	SC_BAD_GATEWAY                     StatusCode = 502
	SC_SERVICE_UNAVAILABLE             StatusCode = 503
	SC_GATEWAY_TIMEOUT                 StatusCode = 504
	SC_HTTP_VERSION_NOT_SUPPORTED      StatusCode = 505
)

var StatusMessageMap = map[StatusCode]string{
	SC_OK:                    "OK",
	SC_ACCEPTED:              "Accepted",
	SC_BAD_GATEWAY:           "Bad Gateway",
	SC_BAD_REQUEST:           "Bad Request",
	SC_CONTINUE:              "Continue",
	SC_FORBIDDEN:             "Forbidden",
	SC_INTERNAL_SERVER_ERROR: "Internal Server Error",
	SC_METHOD_NOT_ALLOWED:    "Method Not Allowed",
	SC_NOT_FOUND:             "Not Found",
	SC_NOT_IMPLEMENTED:       "Not Implemented",
	SC_REQUEST_URI_TOO_LONG:  "Request URI Too Long",
	SC_SERVICE_UNAVAILABLE:   "Service Unavailable",
	SC_UNAUTHORIZED:          "Unauthorized",
}
var DefaultMessage = "HTTP Response Status %d"

type HttpServletResponse interface {
	servlet.ServletResponse

	AddCookie(c *cookie.Cookie)

	SetHeader(name, value string)
	GetHeader(name string) string
	GetHeaderNames() []string

	SetStatus(status int32)
	GetStatus() int32
}
