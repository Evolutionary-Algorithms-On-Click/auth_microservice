package controller

import (
	"evolve/util"
	"evolve/util/auth"
	"net/http" 
	"evolve/modules/team"
)

//to create a new team
func CreateTeam(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	logger.InfoCtx(req, "CreateTeam API called.")

	token, err := req.Cookie("t")

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	payLoad, err := auth.ValidateToken(token.Value)

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}
	
	if payLoad["purpose"] != "login" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	//should change it to data
	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	createTeamReq, err := util.FromJson[team.CreateTeamReq](data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = createTeamReq.CreateTeam(req.Context(), payLoad)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusOK, "Team Creation Successful.", nil)


}

//to get the list of teams created by a user
func GetTeams(res http.ResponseWriter, req *http.Request) {
	logger := util.SharedLogger
	logger.InfoCtx(req, "GetTeams API called.")

	token, err := req.Cookie("t")

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	payLoad, err := auth.ValidateToken(token.Value)

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}

	if payLoad["purpose"] != "login" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	teamsInfo,err := team.GetTeams(req.Context(), payLoad)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	util.JSONResponse(res, http.StatusOK, "Success", teamsInfo)

}

//to get the members and team metadata

func GetTeamMembers(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	logger.InfoCtx(req, "GetTeamMembers API called.")

	token, err := req.Cookie("t")

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	payLoad, err := auth.ValidateToken(token.Value)

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}

	if payLoad["purpose"] != "login" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	getTeamMembersObj, err := util.FromJson[team.GetTeamMembersReq](data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	teamMetadata,err := getTeamMembersObj.GetTeamMembers(req.Context(), payLoad)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	util.JSONResponse(res, http.StatusOK, "Success", teamMetadata)
}

func AddTeamMembers(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	logger.InfoCtx(req, "AddTeamMembers API called.")

	token, err := req.Cookie("t")

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	payLoad, err := auth.ValidateToken(token.Value)

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}

	if payLoad["purpose"] != "login" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	data, err := util.Body(req)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	addTeamMembersReq, err := util.FromJson[team.AddMembersReq](data)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := addTeamMembersReq.AddTeamMembers(req.Context(), payLoad)

	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusOK, result, nil)



}

func DeleteTeamMembers(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	logger.InfoCtx(req, "DeleteTeamMembers API called.")

	token, err := req.Cookie("t")

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	payLoad, err := auth.ValidateToken(token.Value)

	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}

	if payLoad["purpose"] != "login" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	deleteTeamMembersReq, err := util.FromJson[team.DeleteTeamMembersReq](data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := deleteTeamMembersReq.DeleteTeamMembers(req.Context(), payLoad)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusOK, result, nil)
}