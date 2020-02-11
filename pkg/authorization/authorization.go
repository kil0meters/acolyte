package authorization

// PermissionLevel enum for different forum permissions
type PermissionLevel string

const (
	// Admin legends
	Admin PermissionLevel = "AUTH_ADMIN"
	// Moderator chat moderators
	Moderator PermissionLevel = "AUTH_MODERATOR"
	// Standard plebs
	Standard PermissionLevel = "AUTH_STANDARD"
	// LoggedOut supa plebs
	LoggedOut PermissionLevel = "AUTH_LOGGED_OUT"
	// Banned banned plebs
	Banned PermissionLevel = "AUTH_BANNED"
)

// AtLeast tests if a PermissionLevel is at least another PermissionLevel
func (permissions PermissionLevel) AtLeast(minimumPermission PermissionLevel) bool {
	if minimumPermission == Admin {
		return permissions == Admin
	} else if minimumPermission == Moderator {
		return permissions == Moderator || permissions == Admin
	} else if minimumPermission == Standard {
		return permissions == Standard || permissions == Moderator || permissions == Admin
	} else if minimumPermission == LoggedOut {
		return permissions == LoggedOut || permissions == Standard || permissions == Moderator || permissions == Admin
	}
	return true
}
