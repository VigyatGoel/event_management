package models

// Registration represents an attendee registration for an event
type Registration struct {
	ID               int    `json:"id"`
	EventID          int    `json:"eventId"`
	UserID           int    `json:"userId"`
	RegistrationDate string `json:"registrationDate"`
	Status           string `json:"status"` // confirmed, pending, cancelled
}

// RegistrationWithUserDetails includes user information with the registration
type RegistrationWithUserDetails struct {
	Registration
	UserName string `json:"userName"`
	Email    string `json:"email"`
}
