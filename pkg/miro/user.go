package miro

type User struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Active         bool   `json:"active"`
	License        string `json:"license"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	LastActivityAt string `json:"lastActivityAt"`
}
