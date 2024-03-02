package domain

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
	GuestRole Role = "guest"
)

// DefineRoles defines the roles that users can have
type DefineRoles map[Role][]string

// Roles defines the roles map
var Roles = DefineRoles{
	AdminRole: {"admin", "user", "guest"},
	UserRole:  {"user", "guest"},
	GuestRole: {"guest"},
}
