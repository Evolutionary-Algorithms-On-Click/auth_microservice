package routes

const (
	BASE = "/api"
)

const (
	TEST     = BASE + "/test"
	REGISTER = BASE + "/register"
	VERIFY   = REGISTER + "/verify"
	LOGIN    = BASE + "/login"
)

const (
	TEAMS                 = BASE + "/teams"           // POST: create team, GET: list teams
	TEAM_DETAILS          = TEAMS + "/:id"            // GET: get team details + members
	TEAM_MEMBERS          = TEAM_DETAILS + "/members" // POST: add one or more members
	TEAM_MEMBER_OPERATION = TEAM_MEMBERS + "/:userId" // DELETE: remove member
)
