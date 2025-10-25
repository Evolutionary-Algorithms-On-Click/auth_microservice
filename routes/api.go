package routes

const (
	BASE = "/api"
)

const (
	TEST     = BASE + "/test"
	REGISTER = BASE + "/register"
	VERIFY   = REGISTER + "/verify"
	LOGIN    = BASE + "/login"

	// Password Reset Routes
	PASSWORD        = BASE + "/password"
	PASSWORD_RESET  = PASSWORD + "/reset"
	PASSWORD_VERIFY = PASSWORD_RESET + "/verify"
)
