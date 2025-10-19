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
	RESET_REQUEST = BASE + "/reset-password-request"
	RESET_VERIFY  = BASE + "/reset-password-verify"
)
