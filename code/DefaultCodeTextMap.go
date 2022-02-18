package code


const UndefinedCodeText = "undefined code"

var statusText= map[string]map[Code]string{
	"en": statusTextEN,
}


// RFC code
const (
	StatusContinue           Code = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols Code = 101 // RFC 7231, 6.2.2
	StatusProcessing         Code = 102 // RFC 2518, 10.1
	StatusEarlyHints         Code = 103 // RFC 8297

	OK                         Code = 200 // RFC 7231, 6.3.1
	Created                    Code = 201 // RFC 7231, 6.3.2
	Accepted                   Code = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo Code = 203 // RFC 7231, 6.3.4
	NoContent                  Code = 204 // RFC 7231, 6.3.5
	StatusResetContent         Code = 205 // RFC 7231, 6.3.6
	StatusPartialContent       Code = 206 // RFC 7233, 4.1
	StatusMultiStatus          Code = 207 // RFC 4918, 11.1
	StatusAlreadyReported      Code = 208 // RFC 5842, 7.1
	StatusIMUsed               Code = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   Code = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  Code = 301 // RFC 7231, 6.4.2
	StatusFound             Code = 302 // RFC 7231, 6.4.3
	StatusSeeOther          Code = 303 // RFC 7231, 6.4.4
	StatusNotModified       Code = 304 // RFC 7232, 4.1
	StatusUseProxy          Code = 305 // RFC 7231, 6.4.5
	StatusTemporaryRedirect Code = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect Code = 308 // RFC 7538, 3

	BadRequest Code = 400 // RFC 7231, 6.5.1
	// 未豋入之類
	Unauthorized          Code = 401 // RFC 7235, 3.1
	StatusPaymentRequired Code = 402 // RFC 7231, 6.5.2
	// 已豋入但權限不族
	Forbidden                          Code = 403 // RFC 7231, 6.5.3
	NotFound                           Code = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             Code = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                Code = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            Code = 407 // RFC 7235, 3.2
	StatusRequestTimeout               Code = 408 // RFC 7231, 6.5.7
	StatusConflict                     Code = 409 // RFC 7231, 6.5.8
	StatusGone                         Code = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               Code = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           Code = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        Code = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            Code = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         Code = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable Code = 416 // RFC 7233, 4.4
	StatusExpectationFailed            Code = 417 // RFC 7231, 6.5.14
	StatusTeapot                       Code = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           Code = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          Code = 422 // RFC 4918, 11.2
	StatusLocked                       Code = 423 // RFC 4918, 11.3
	StatusFailedDependency             Code = 424 // RFC 4918, 11.4
	StatusTooEarly                     Code = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              Code = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         Code = 428 // RFC 6585, 3
	StatusTooManyRequests              Code = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  Code = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   Code = 451 // RFC 7725, 3

	ServerError                   Code = 500 // RFC 7231, 6.6.1
	NotImplemented                Code = 501 // RFC 7231, 6.6.2
	StatusBadGateway              Code = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable      Code = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout          Code = 504 // RFC 7231, 6.6.5
	HTTPVersionNotSupported       Code = 505 // RFC 7231, 6.6.6
	VariantAlsoNegotiates         Code = 506 // RFC 2295, 8.1
	InsufficientStorage           Code = 507 // RFC 4918, 11.5
	LoopDetected                  Code = 508 // RFC 5842, 7.2
	NotExtended                   Code = 510 // RFC 2774, 7
	NetworkAuthenticationRequired Code = 511 // RFC 6585, 6
)
var statusTextEN = map[Code]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "Early Hints",

	OK:                         "OK",
	Created:                    "Created",
	Accepted:                   "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	NoContent:                  "",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",
	StatusMultiStatus:          "Multi-Status",
	StatusAlreadyReported:      "Already Reported",
	StatusIMUsed:               "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	BadRequest:                         "Bad Request",
	Unauthorized:                       "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	Forbidden:                          "Forbidden",
	NotFound:                           "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusMisdirectedRequest:           "Misdirected Request",
	StatusUnprocessableEntity:          "Unprocessable Entity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "Failed Dependency",
	StatusTooEarly:                     "Too Early",
	StatusUpgradeRequired:              "Upgrade Required",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	ServerError:                   "Internal Server Error",
	NotImplemented:                "Not Implemented",
	StatusBadGateway:              "Bad Gateway",
	StatusServiceUnavailable:      "Service Unavailable",
	StatusGatewayTimeout:          "Gateway Timeout",
	HTTPVersionNotSupported:       "HTTP Version Not Supported",
	VariantAlsoNegotiates:         "Variant Also Negotiates",
	InsufficientStorage:           "Insufficient Storage",
	LoopDetected:                  "Loop Detected",
	NotExtended:                   "Not Extended",
	NetworkAuthenticationRequired: "Network Authentication Required",
}

