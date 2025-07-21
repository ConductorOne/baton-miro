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

type ScimUserRole struct {
	Value   string `json:"value"`
	Display string `json:"display"`
	Type    string `json:"type"`
	Primary bool   `json:"primary"`
}

type ScimUserEmail struct {
	Value   string `json:"value"`
	Display string `json:"display"`
	Primary bool   `json:"primary"`
}

type ScimUserGroup struct {
	Value   string `json:"value"`
	Display string `json:"display"`
}

type ScimUserName struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type ScimUser struct {
	Schemas     []string        `json:"schemas"`
	Id          string          `json:"id"`
	UserName    string          `json:"userName"`
	Name        ScimUserName    `json:"name"`
	DisplayName string          `json:"displayName"`
	Active      bool            `json:"active"`
	UserType    string          `json:"userType"`
	Emails      []ScimUserEmail `json:"emails"`
	Groups      []ScimUserGroup `json:"groups"`
	Roles       []ScimUserRole  `json:"roles"`
}

type PatchOp struct {
	Schemas    []string      `json:"schemas"`
	Operations []PatchOpItem `json:"Operations"`
}

type PatchOpItem struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}
