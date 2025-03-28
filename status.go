package webutil

import "net/http"

// Pre-defined HTTP status error constants.
// These provide HTTP status codes as HTTPError objects that
// can be used directly as errors or as http.Handler instances.
const (
	// 4xx Client Errors
	StatusBadRequest                   = HTTPError(http.StatusBadRequest)
	StatusUnauthorized                 = HTTPError(http.StatusUnauthorized)
	StatusPaymentRequired              = HTTPError(http.StatusPaymentRequired)
	StatusForbidden                    = HTTPError(http.StatusForbidden)
	StatusNotFound                     = HTTPError(http.StatusNotFound)
	StatusMethodNotAllowed             = HTTPError(http.StatusMethodNotAllowed)
	StatusNotAcceptable                = HTTPError(http.StatusNotAcceptable)
	StatusProxyAuthRequired            = HTTPError(http.StatusProxyAuthRequired)
	StatusRequestTimeout               = HTTPError(http.StatusRequestTimeout)
	StatusConflict                     = HTTPError(http.StatusConflict)
	StatusGone                         = HTTPError(http.StatusGone)
	StatusLengthRequired               = HTTPError(http.StatusLengthRequired)
	StatusPreconditionFailed           = HTTPError(http.StatusPreconditionFailed)
	StatusRequestEntityTooLarge        = HTTPError(http.StatusRequestEntityTooLarge)
	StatusRequestURITooLong            = HTTPError(http.StatusRequestURITooLong)
	StatusUnsupportedMediaType         = HTTPError(http.StatusUnsupportedMediaType)
	StatusRequestedRangeNotSatisfiable = HTTPError(http.StatusRequestedRangeNotSatisfiable)
	StatusExpectationFailed            = HTTPError(http.StatusExpectationFailed)
	StatusTeapot                       = HTTPError(http.StatusTeapot)
	StatusMisdirectedRequest           = HTTPError(http.StatusMisdirectedRequest)
	StatusUnprocessableEntity          = HTTPError(http.StatusUnprocessableEntity)
	StatusLocked                       = HTTPError(http.StatusLocked)
	StatusFailedDependency             = HTTPError(http.StatusFailedDependency)
	StatusTooEarly                     = HTTPError(http.StatusTooEarly)
	StatusUpgradeRequired              = HTTPError(http.StatusUpgradeRequired)
	StatusPreconditionRequired         = HTTPError(http.StatusPreconditionRequired)
	StatusTooManyRequests              = HTTPError(http.StatusTooManyRequests)
	StatusRequestHeaderFieldsTooLarge  = HTTPError(http.StatusRequestHeaderFieldsTooLarge)
	StatusUnavailableForLegalReasons   = HTTPError(http.StatusUnavailableForLegalReasons)

	// 5xx Server Errors
	StatusInternalServerError           = HTTPError(http.StatusInternalServerError)
	StatusNotImplemented                = HTTPError(http.StatusNotImplemented)
	StatusBadGateway                    = HTTPError(http.StatusBadGateway)
	StatusServiceUnavailable            = HTTPError(http.StatusServiceUnavailable)
	StatusGatewayTimeout                = HTTPError(http.StatusGatewayTimeout)
	StatusHTTPVersionNotSupported       = HTTPError(http.StatusHTTPVersionNotSupported)
	StatusVariantAlsoNegotiates         = HTTPError(http.StatusVariantAlsoNegotiates)
	StatusInsufficientStorage           = HTTPError(http.StatusInsufficientStorage)
	StatusLoopDetected                  = HTTPError(http.StatusLoopDetected)
	StatusNotExtended                   = HTTPError(http.StatusNotExtended)
	StatusNetworkAuthenticationRequired = HTTPError(http.StatusNetworkAuthenticationRequired)
)
