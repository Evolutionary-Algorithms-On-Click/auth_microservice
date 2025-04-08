package controller

import (
	"evolve/modules"
	"evolve/util"
	"evolve/util/auth"
	"net/http"
	"strings"
)

func CreateTeam(res http.ResponseWriter, req *http.Request) {
	var logger = util.NewLogger()
	logger.Info("CreateTeam API called.")

	// Only POST method is allowed
	if req.Method != http.MethodPost {
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get token from cookie
	cookie, err := req.Cookie("t")
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "No token provided", nil)
		return
	}

	// Validate token
	userClaims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userID := userClaims["id"]
	if userID == "" {
		util.JSONResponse(res, http.StatusUnauthorized, "Invalid token payload", nil)
		return
	}

	// Parse request body
	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	teamReq, err := modules.TeamReqFromJSON(data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Create team
	team, err := teamReq.CreateTeam(req.Context(), userID)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusCreated, "Team created successfully", team)
}

func ListTeams(res http.ResponseWriter, req *http.Request) {
	var logger = util.NewLogger()
	logger.Info("ListTeams API called.")

	// Only GET method is allowed
	if req.Method != http.MethodGet {
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get token from cookie
	cookie, err := req.Cookie("t")
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "No token provided", nil)
		return
	}

	// Validate token
	userClaims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userID := userClaims["id"]
	if userID == "" {
		util.JSONResponse(res, http.StatusUnauthorized, "Invalid token payload", nil)
		return
	}

	// Get teams by user
	teams, err := modules.GetTeamsByUser(req.Context(), userID)
	if err != nil {
		util.JSONResponse(res, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusOK, "Teams retrieved successfully", teams)
}

// Teams handles both creating teams and listing teams
func Teams(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	// Handle sub-routes if they exist (like /:id/members)
	if strings.Contains(path, "/members") || strings.Contains(path, "/") && len(strings.Split(path, "/")) > 4 {
		util.JSONResponse(res, http.StatusNotFound, "Route not implemented", nil)
		return
	}

	switch req.Method {
	case http.MethodPost:
		CreateTeam(res, req)
	case http.MethodGet:
		ListTeams(res, req)
	default:
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}