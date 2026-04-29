package domain

type EventPersonnel struct {
	ID		string	`json:"id"`
	EventID		string	`json:"event_id"`
	Name		string	`json:"name"`
	Role		string	`json:"role"`
	Image		string	`json:"image"`
	ProfileLink	string	`json:"profile_link"`
}
