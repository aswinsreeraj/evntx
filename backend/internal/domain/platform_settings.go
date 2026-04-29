package domain

import "time"

const PlatformSettingsID = "00000000-0000-0000-0000-000000000002"

type PlatformFeeType string

const (
	PlatformFeeTypeFixed		PlatformFeeType	= "fixed"
	PlatformFeeTypePercentage	PlatformFeeType	= "percentage"
)

type PlatformSettings struct {
	ID					string		`json:"id"`
	EnableUserRegistration			bool		`json:"enable_user_registration"`
	AllowGoogleLogin			bool		`json:"allow_google_login"`
	RequireAdminApprovalForOrganizers	bool		`json:"require_admin_approval_for_organizers"`
	RequireAdminApprovalForEvents		bool		`json:"require_admin_approval_for_events"`
	RefundWindowDays			int		`json:"refund_window_days"`
	AllowEventCancellation			bool		`json:"allow_event_cancellation"`
	PlatformFeeValue			float64		`json:"platform_fee_value"`
	PlatformFeeType				PlatformFeeType	`json:"platform_fee_type"`
	UpdatedAt				time.Time	`json:"updated_at"`
}

type PaymentSettings struct {
	ID		string			`json:"id"`
	Provider	string			`json:"provider"`
	IsEnabled	bool			`json:"is_enabled"`
	Config		map[string]interface{}	`json:"config"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
}
