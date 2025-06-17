package constants

var (
	JSONResponse    = "JSONResponse"
	Db              = "Database"
	CustomValidator = "CustomValidator"
	Bearer          = "Bearer"
	VerifyToken     = "VerifyToken"
	UserService     = "UserService"
	AuthService     = "AuthService"
)

const (
	RoleAdmin = iota
	RoleDirector
	RoleManager
	RoleStaff
)
