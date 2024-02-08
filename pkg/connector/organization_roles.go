package connector

const (
	internalAdminOrganizationRole         = "organization_internal_admin"
	internalUserOrganizationRole          = "organization_internal_user"
	internalExternalUserOrganizationRole  = "organization_external_user"
	internalTeamGuestUserOrganizationRole = "organization_team_guest_user"
	unknownOrganizationRole               = "unknown"
)

var organizationRoles = []string{
	internalAdminOrganizationRole,
	internalUserOrganizationRole,
	internalExternalUserOrganizationRole,
	internalTeamGuestUserOrganizationRole,
	unknownOrganizationRole,
}
