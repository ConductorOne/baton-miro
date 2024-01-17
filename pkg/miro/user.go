package miro

type User struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Active         bool   `json:"active"`
	License        string `json:"license"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	LastActivityAt string `json:"lastActivityAt"`
}
