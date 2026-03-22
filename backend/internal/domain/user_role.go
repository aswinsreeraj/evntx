package domain

type UserRole string

const (
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)
