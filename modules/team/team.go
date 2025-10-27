package team

import (
	"context"
	"encoding/json"
	"fmt"
	"evolve/db/connection"
	"evolve/util"
)

type CreateTeamReq struct {
	TeamName string `json:"teamName"`
	TeamDesc    string `json:"teamDesc"`
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
	Email string 	`json:"email"`
}

type TeamData struct {
	TeamId string `json:"teamId"`
	TeamName string `json:"teamName"`
	TeamDesc string `json:"teamDesc"`
}


func CreateTeamReqFromJson(data map[string]any) (*CreateTeamReq, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var createTeamReq *CreateTeamReq
	if err := json.Unmarshal(jsonData, &createTeamReq); err != nil {
		return nil, err
	}

	return createTeamReq, nil
}


func (c *CreateTeamReq) CreateTeam(ctx context.Context, user map[string]string) error {
	
	logger := util.SharedLogger

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("CreateTeam: failed to get pool connection: %v", err), err)
		return fmt.Errorf("something went wrong")
	}

	

	//inserting team info in team table
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
		return nil,fmt.Errorf("something went wrong")
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
		return nil,fmt.Errorf("something went wrong")
	}

	rows, err := db.Query(ctx, "SELECT T.TEAMID, T.TEAMNAME, T.TEAMDESC, M.MEMBERID, U.USERNAME,U.EMAIL, M.ROLE FROM TEAM T JOIN TEAMMEMBERS M ON T.TEAMID = M.TEAMID JOIN USERS U ON M.MEMBERID = U.ID WHERE T.CREATEDBY = $1 AND T.TEAMNAME = $2", payLoad["id"], g.TeamName)

	if err != nil {
		logger.Error(fmt.Sprintf("GetTeamMembers: failed to get teams: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
		}
	
	var result map[string]any

	var membersInfo []TeamMembers
	var teamMetadata TeamData

	fmt.Println(rows)
	for rows.Next() {
		var teamData TeamMembers
		
		//getting team members list
		err := rows.Scan(&teamMetadata.TeamId, &teamMetadata.TeamName, &teamMetadata.TeamDesc ,&teamData.UserName, &teamData.Email)
		if err != nil {
            logger.Error(fmt.Sprintf("GetTeamMembersReq: failed to get team data: %v", err), err)
            return nil, fmt.Errorf("something went wrong")
        }

		membersInfo = append(membersInfo, teamData)
	}

	// result1, err := json.Marshal(membersInfo)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("GetTeamMembers: failed to convert MembersInfo to json: %v", err), err)
	// 	return nil, fmt.Errorf("something went wrong")
	// }

	// result2, errr := json.Marshal(teamMetadata)

	
	result["members"] = membersInfo
	result["teamData"] = teamMetadata

	return result, nil
	





}