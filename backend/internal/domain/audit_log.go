package domain

import "time"

type AuditActionTag string

const (
	ActionTagEvent     AuditActionTag = "EVENT"
	ActionTagUser      AuditActionTag = "USER"
	ActionTagOrganizer AuditActionTag = "ORGANIZER"
	ActionTagPayout    AuditActionTag = "PAYOUT"
	ActionTagSettings  AuditActionTag = "SETTINGS"
)

type AuditLog struct {
	ID        string         `json:"id"`
	AdminID   string         `json:"admin_id"`
	AdminName string         `json:"admin_name"`
	Action    string         `json:"action"`
	ActionTag AuditActionTag `json:"action_tag"`
	Details   string         `json:"details"`
	IPAddress string         `json:"ip_address"`
	Timestamp time.Time      `json:"timestamp"`
}
