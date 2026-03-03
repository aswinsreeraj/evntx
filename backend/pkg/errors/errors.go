package errors

const (
	InvalidRequestBody     = "EVT_001"
	ResourceNotFound       = "EVT_002"
	UnauthorizedAccess     = "EVT_003"
	ForbiddenAction        = "EVT_004"
	DuplicateResource      = "EVT_005"
	InvalidStateTransition = "EVT_006"
	PaymentFailed          = "EVT_007"
	InsufficientBalance    = "EVT_008"
	TicketSoldOut          = "EVT_009"
	SessionExpired         = "EVT_010"
	InvalidOTP             = "EVT_011"
	EventNotLive           = "EVT_012"
	RateLimitExceeded      = "EVT_013"
)
