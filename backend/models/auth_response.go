package models

type AuthResponse struct {
	Message string `json:"message"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Token   string `json:"token,omitempty"`
}