package errorx

import (
	"gin-admin/locales"
	"net/http"

	"github.com/epkgs/i18n/errors"
)

func defineHttpError(httpStatus int, defaultMessage any, errcode ...int) errors.Error {

	var code int
	if len(errcode) > 0 {
		code = errcode[0]
	}
	if code <= 0 {
		code = httpStatus
	}

	if code == 200 {
		code = 1
	}

	err := errors.New(defaultMessage)

	return errors.WithHttpStatus(
		errors.WithCode(err, code),
		httpStatus,
	)
}

var (
	ErrContinue           = defineHttpError(http.StatusContinue, locales.Http.Str("Continue"))
	ErrSwitchingProtocols = defineHttpError(http.StatusSwitchingProtocols, locales.Http.Str("Switching Protocols"))
	ErrProcessing         = defineHttpError(http.StatusProcessing, locales.Http.Str("Processing"))
	ErrEarlyHints         = defineHttpError(http.StatusEarlyHints, locales.Http.Str("Early Hints"))

	ErrOK                   = defineHttpError(http.StatusOK, locales.Http.Str("OK"))
	ErrCreated              = defineHttpError(http.StatusCreated, locales.Http.Str("Created"))
	ErrAccepted             = defineHttpError(http.StatusAccepted, locales.Http.Str("Accepted"))
	ErrNonAuthoritativeInfo = defineHttpError(http.StatusNonAuthoritativeInfo, locales.Http.Str("Non-Authoritative Information"))
	ErrNoContent            = defineHttpError(http.StatusNoContent, locales.Http.Str("No Content"))
	ErrResetContent         = defineHttpError(http.StatusResetContent, locales.Http.Str("Reset Content"))
	ErrPartialContent       = defineHttpError(http.StatusPartialContent, locales.Http.Str("Partial Content"))
	ErrMultiStatus          = defineHttpError(http.StatusMultiStatus, locales.Http.Str("Multi-Status"))
	ErrAlreadyReported      = defineHttpError(http.StatusAlreadyReported, locales.Http.Str("Already Reported"))
	ErrIMUsed               = defineHttpError(http.StatusIMUsed, locales.Http.Str("IM Used"))

	ErrMultipleChoices   = defineHttpError(http.StatusMultipleChoices, locales.Http.Str("Multiple Choices"))
	ErrMovedPermanently  = defineHttpError(http.StatusMovedPermanently, locales.Http.Str("Moved Permanently"))
	ErrFound             = defineHttpError(http.StatusFound, locales.Http.Str("Found"))
	ErrSeeOther          = defineHttpError(http.StatusSeeOther, locales.Http.Str("See Other"))
	ErrNotModified       = defineHttpError(http.StatusNotModified, locales.Http.Str("Not Modified"))
	ErrUseProxy          = defineHttpError(http.StatusUseProxy, locales.Http.Str("Use Proxy"))
	ErrTemporaryRedirect = defineHttpError(http.StatusTemporaryRedirect, locales.Http.Str("Temporary Redirect"))
	ErrPermanentRedirect = defineHttpError(http.StatusPermanentRedirect, locales.Http.Str("Permanent Redirect"))

	ErrBadRequest                   = defineHttpError(http.StatusBadRequest, locales.Http.Str("Bad Request"))
	ErrUnauthorized                 = defineHttpError(http.StatusUnauthorized, locales.Http.Str("Unauthorized"))
	ErrPaymentRequired              = defineHttpError(http.StatusPaymentRequired, locales.Http.Str("Payment Required"))
	ErrForbidden                    = defineHttpError(http.StatusForbidden, locales.Http.Str("Forbidden"))
	ErrNotFound                     = defineHttpError(http.StatusNotFound, locales.Http.Str("Not Found"))
	ErrMethodNotAllowed             = defineHttpError(http.StatusMethodNotAllowed, locales.Http.Str("Method Not Allowed"))
	ErrNotAcceptable                = defineHttpError(http.StatusNotAcceptable, locales.Http.Str("Not Acceptable"))
	ErrProxyAuthRequired            = defineHttpError(http.StatusProxyAuthRequired, locales.Http.Str("Proxy Authentication Required"))
	ErrRequestTimeout               = defineHttpError(http.StatusRequestTimeout, locales.Http.Str("Request Timeout"))
	ErrConflict                     = defineHttpError(http.StatusConflict, locales.Http.Str("Conflict"))
	ErrGone                         = defineHttpError(http.StatusGone, locales.Http.Str("Gone"))
	ErrLengthRequired               = defineHttpError(http.StatusLengthRequired, locales.Http.Str("Length Required"))
	ErrPreconditionFailed           = defineHttpError(http.StatusPreconditionFailed, locales.Http.Str("Precondition Failed"))
	ErrRequestEntityTooLarge        = defineHttpError(http.StatusRequestEntityTooLarge, locales.Http.Str("Request Entity Too Large"))
	ErrRequestURITooLong            = defineHttpError(http.StatusRequestURITooLong, locales.Http.Str("Request URI Too Long"))
	ErrUnsupportedMediaType         = defineHttpError(http.StatusUnsupportedMediaType, locales.Http.Str("Unsupported Media Type"))
	ErrRequestedRangeNotSatisfiable = defineHttpError(http.StatusRequestedRangeNotSatisfiable, locales.Http.Str("Requested Range Not Satisfiable"))
	ErrExpectationFailed            = defineHttpError(http.StatusExpectationFailed, locales.Http.Str("Expectation Failed"))
	ErrTeapot                       = defineHttpError(http.StatusTeapot, locales.Http.Str("I'm a teapot"))
	ErrMisdirectedRequest           = defineHttpError(http.StatusMisdirectedRequest, locales.Http.Str("Misdirected Request"))
	ErrUnprocessableEntity          = defineHttpError(http.StatusUnprocessableEntity, locales.Http.Str("Unprocessable Entity"))
	ErrLocked                       = defineHttpError(http.StatusLocked, locales.Http.Str("Locked"))
	ErrFailedDependency             = defineHttpError(http.StatusFailedDependency, locales.Http.Str("Failed Dependency"))
	ErrTooEarly                     = defineHttpError(http.StatusTooEarly, locales.Http.Str("Too Early"))
	ErrUpgradeRequired              = defineHttpError(http.StatusUpgradeRequired, locales.Http.Str("Upgrade Required"))
	ErrPreconditionRequired         = defineHttpError(http.StatusPreconditionRequired, locales.Http.Str("Precondition Required"))
	ErrTooManyRequests              = defineHttpError(http.StatusTooManyRequests, locales.Http.Str("Too Many Requests"))
	ErrRequestHeaderFieldsTooLarge  = defineHttpError(http.StatusRequestHeaderFieldsTooLarge, locales.Http.Str("Request Header Fields Too Large"))
	ErrUnavailableForLegalReasons   = defineHttpError(http.StatusUnavailableForLegalReasons, locales.Http.Str("Unavailable For Legal Reasons"))

	ErrInternalServerError           = defineHttpError(http.StatusInternalServerError, locales.Http.Str("Internal Server Error"))
	ErrNotImplemented                = defineHttpError(http.StatusNotImplemented, locales.Http.Str("Not Implemented"))
	ErrBadGateway                    = defineHttpError(http.StatusBadGateway, locales.Http.Str("Bad Gateway"))
	ErrServiceUnavailable            = defineHttpError(http.StatusServiceUnavailable, locales.Http.Str("Service Unavailable"))
	ErrGatewayTimeout                = defineHttpError(http.StatusGatewayTimeout, locales.Http.Str("Gateway Timeout"))
	ErrHTTPVersionNotSupported       = defineHttpError(http.StatusHTTPVersionNotSupported, locales.Http.Str("HTTP Version Not Supported"))
	ErrVariantAlsoNegotiates         = defineHttpError(http.StatusVariantAlsoNegotiates, locales.Http.Str("Variant Also Negotiates"))
	ErrInsufficientStorage           = defineHttpError(http.StatusInsufficientStorage, locales.Http.Str("Insufficient Storage"))
	ErrLoopDetected                  = defineHttpError(http.StatusLoopDetected, locales.Http.Str("Loop Detected"))
	ErrNotExtended                   = defineHttpError(http.StatusNotExtended, locales.Http.Str("Not Extended"))
	ErrNetworkAuthenticationRequired = defineHttpError(http.StatusNetworkAuthenticationRequired, locales.Http.Str("Network Authentication Required"))
)
