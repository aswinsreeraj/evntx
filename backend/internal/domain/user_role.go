package domain

type UserRole string

const (
	RoleGoer      UserRole = "goer"
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)
