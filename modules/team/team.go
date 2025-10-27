package team

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"fmt"
)

type CreateTeamReq struct {
	TeamName string `json:"teamName"`
	TeamDesc string `json:"teamDesc"`
}

type TeamInfo struct {
	TeamID      string `json:"teamId"`
	TeamDesc    string `json:"teamDesc"`
	MemberCount int    `json:"memberCount"`
}

type GetTeamMembersReq struct {
	TeamName string `json:"teamName"`
}

type TeamMembers struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type TeamData struct {
	TeamId   string `json:"teamId"`
	TeamName string `json:"teamName"`
	TeamDesc string `json:"teamDesc"`
}

type DeleteTeamMembersReq struct {
	TeamName    string   `json:"teamName"`
	TeamMembers []string `json:"teamMembers"` // a list of usernames
}

type AddMembersReq struct {
	TeamName    string   `json:"teamName"`
	TeamMembers []string `json:"teamMembers"`
}

func (c *CreateTeamReq) CreateTeam(ctx context.Context, user map[string]string) error {

	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("CreateTeam: failed to get pool connection: %v", err), err)
		return fmt.Errorf("something went wrong")
	}

	var teamID string
	err = db.QueryRow(ctx, "INSERT INTO team (teamName, teamDesc, createdBy) VALUES ($1, $2, $3) RETURNING teamID", c.TeamName, c.TeamDesc, user["id"]).Scan(&teamID)

	if err != nil {
		logger.Error(fmt.Sprintf("Createteam: failed to create Team: %v", err), err)
		return fmt.Errorf("something went wrong")
	}

	//inserting the admin into teamMembers table
	_, err = db.Exec(ctx, "INSERT INTO teamMembers (memberId, teamID, role) VALUES ($1, $2, $3)", user["id"], teamID, "Admin")

	if err != nil {
		logger.Error(fmt.Sprintf("Createteam: failed to Insert into teamMembers Table: %v", err), err)
		return fmt.Errorf("something went wrong")
	}
	return nil

}

func GetTeams(ctx context.Context, user map[string]string) ([]map[string]any, error) {
	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("GetTeams: failed to get pool connection: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	rows, err := db.Query(ctx, "SELECT T.TEAMID, T.TEAMDESC, COUNT(*) OVER (PARTITION BY M.TEAMID) FROM TEAM T JOIN TEAMMEMBERS M ON T.TEAMID = M.TEAMID AND T.CREATEDBY = $1", user["id"])

	if err != nil {
		logger.Error(fmt.Sprintf("GetTeams: failed to get teams: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	var teams []TeamInfo

	for rows.Next() {
		var team TeamInfo
		err := rows.Scan(&team.TeamID, &team.TeamDesc, &team.MemberCount)
		if err != nil {
			logger.Error(fmt.Sprintf("GetTeams: failed to get teams: %v", err), err)
			return nil, fmt.Errorf("something went wrong")
		}

		teams = append(teams, team)
	}

	result, err := json.Marshal(teams)
	if err != nil {
		logger.Error(fmt.Sprintf("GetTeams: failed to convert TeamInfo to json: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	var teamMap []map[string]any

	err = json.Unmarshal(result, &teamMap)

	return teamMap, nil
}

func (g *GetTeamMembersReq) GetTeamMembers(ctx context.Context, payLoad map[string]string) (map[string]any, error) {

	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("GetTeamMembers: failed to get pool connection: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	rows, err := db.Query(ctx, "SELECT T.TEAMID, T.TEAMNAME, T.TEAMDESC, M.MEMBERID, U.USERNAME, U.EMAIL FROM TEAM T JOIN TEAMMEMBERS M ON T.TEAMID = M.TEAMID JOIN USERS U ON M.MEMBERID = U.ID WHERE T.CREATEDBY = $1 AND T.TEAMNAME = $2", payLoad["id"], g.TeamName)

	if err != nil {
		logger.Error(fmt.Sprintf("GetTeamMembers: failed to get teams: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	result := make(map[string]any)

	var membersInfo []TeamMembers
	var teamMetadata TeamData

	firstRow := true
	for rows.Next() {
		var teamData TeamMembers
		var teamID, memberID string

		//getting team members list
		err := rows.Scan(&teamID, &teamMetadata.TeamName, &teamMetadata.TeamDesc, &memberID, &teamData.UserName, &teamData.Email)
		if err != nil {
			logger.Error(fmt.Sprintf("GetTeamMembersReq: failed to get team data: %v", err), err)
			return nil, fmt.Errorf("something went wrong")
		}

		if firstRow {
			teamMetadata.TeamId = teamID
			firstRow = false
		}

		membersInfo = append(membersInfo, teamData)
	}

	result["members"] = membersInfo
	result["teamData"] = teamMetadata

	return result, nil

}

func (a *AddMembersReq) AddTeamMembers(ctx context.Context, payLoad map[string]string) (string, error) {
	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("AddTeamMembers: failed to get pool connection: %v", err), err)
		return "", fmt.Errorf("something went wrong")
	}

	var teamID string
	err = db.QueryRow(ctx, "SELECT teamID FROM team WHERE teamName = $1 AND createdBy = $2", a.TeamName, payLoad["id"]).Scan(&teamID)
	if err != nil {
		logger.Error(fmt.Sprintf("AddTeamMembers: Invalid Team ID: %v", err), err)
		return "", fmt.Errorf("Team Not Found")
	}

	// Add each member to the team
	addedCount := 0
	for _, member := range a.TeamMembers {

		var userID string
		err = db.QueryRow(ctx, "SELECT id FROM users WHERE userName = $1", member).Scan(&userID)
		if err != nil {
			logger.Error(fmt.Sprintf("AddTeamMembers: User not found: %v", member), err)
			continue
		}

		_, err = db.Exec(ctx, "INSERT INTO teamMembers (memberId, teamID, role) VALUES ($1, $2, $3)", userID, teamID, "Member")
		if err != nil {
			logger.Error(fmt.Sprintf("AddTeamMembers: failed to add member %v: %v", member, err), err)
			continue
		}
		addedCount++
	}

	return fmt.Sprintf("Added %d members into the team", addedCount), nil
}

func (d *DeleteTeamMembersReq) DeleteTeamMembers(ctx context.Context, payLoad map[string]string) (string, error) {

	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("DeleteTeamMembers: failed to get pool connection: %v", err), err)
		return "", fmt.Errorf("something went wrong")
	}

	var teamID string
	err = db.QueryRow(ctx, "SELECT teamID FROM team WHERE teamName = $1 AND createdBy = $2", d.TeamName, payLoad["id"]).Scan(&teamID)
	if err != nil {
		logger.Error(fmt.Sprintf("DeleteTeamMembers: Invalid Team ID: %v", err), err)
		return "", fmt.Errorf("Team Not Found")
	}

	result, err := db.Exec(ctx, "DELETE FROM teamMembers WHERE teamID = $1 AND memberId IN (SELECT id FROM users WHERE userName = ANY($2))", teamID, d.TeamMembers)
	if err != nil {
		logger.Error(fmt.Sprintf("DeleteTeamMembers: failed to delete members: %v", err), err)
		return "", fmt.Errorf("something went wrong")
	}

	deletedCount := result.RowsAffected()
	return fmt.Sprintf("Deleted %d members from team", deletedCount), nil
}
