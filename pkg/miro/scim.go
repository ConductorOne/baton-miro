package miro

import "fmt"

// CreateUserRequest defines the payload for creating a new user via SCIM.
type CreateUserRequest struct {
	Schemas  []string    `json:"schemas"`
	UserName string      `json:"userName"`
	Name     RequestName `json:"name"`
}

// RequestName defines the name of a user.
type RequestName struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

// SCIMError defines the error response from the SCIM API.
type SCIMError struct {
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
	Detail  string   `json:"detail"`
}

// Error returns the error message for a SCIM error.
func (e *SCIMError) Error() string {
	return fmt.Sprintf("miro scim error %s: %s", e.Status, e.Detail)
}

// Message returns the detail message for a SCIM error.
func (e *SCIMError) Message() string {
	return e.Detail
}
