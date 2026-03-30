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
	BookingExpired         = "EVT_014"
	InternalServerError    = "EVT_999"
)

var (
	ErrInvalidRequestBody     = New(400, InvalidRequestBody, "Invalid request body")
	ErrResourceNotFound       = New(404, ResourceNotFound, "Resource not found")
	ErrUnauthorizedAccess     = New(401, UnauthorizedAccess, "Unauthorized access")
	ErrForbiddenAction        = New(403, ForbiddenAction, "Forbidden action")
	ErrDuplicateResource      = New(409, DuplicateResource, "Duplicate resource")
	ErrInvalidStateTransition = New(400, InvalidStateTransition, "Invalid state transition")
	ErrPaymentFailed          = New(402, PaymentFailed, "Payment failed")
	ErrInsufficientBalance    = New(402, InsufficientBalance, "Insufficient balance")
	ErrTicketSoldOut          = New(400, TicketSoldOut, "Ticket sold out")
	ErrSessionExpired         = New(401, SessionExpired, "Session expired")
	ErrInvalidOTP             = New(400, InvalidOTP, "Invalid OTP")
	ErrEventNotLive           = New(400, EventNotLive, "Event not live")
	ErrRateLimitExceeded      = New(429, RateLimitExceeded, "Rate limit exceeded")
	ErrBookingExpired         = New(400, BookingExpired, "Booking has expired")
	ErrInternalServerError    = New(500, InternalServerError, "Internal server error")
)
