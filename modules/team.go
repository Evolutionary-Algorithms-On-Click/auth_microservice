package modules

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"fmt"
	"time"
)

//Struct for a Team
type Team struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type CreateTeamReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func TeamReqFromJSON(data map[string]any) (*CreateTeamReq, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var teamReq *CreateTeamReq
	if err := json.Unmarshal(jsonData, &teamReq); err != nil {
		return nil, err
	}
	return teamReq, nil
}

func (t *CreateTeamReq) validate() error {
	if len(t.Name) == 0 {
		return fmt.Errorf("team name is required")
	}
	return nil
}

func (t *CreateTeamReq) CreateTeam(ctx context.Context, userID string) (*Team, error) {
	var logger = util.NewLogger()
	if err := t.validate(); err != nil {
		return nil, err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("CreateTeam: failed to get pool connection: %v", err))
		return nil, fmt.Errorf("something went wrong")
	}

	var team Team
	err = db.QueryRow(ctx,
		"INSERT INTO teams (name, description, createdBy) VALUES ($1, $2, $3) RETURNING id, name, description, createdBy, createdAt, updatedAt",
		t.Name, t.Description, userID).Scan(
		&team.ID, &team.Name, &team.Description, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		logger.Error(fmt.Sprintf("CreateTeam: failed to insert team: %v", err))
		return nil, fmt.Errorf("failed to create team")
	}

	// Add the creator as a team member
	_, err = db.Exec(ctx,
		"INSERT INTO team_members (teamID, userID) VALUES ($1, $2)",
		team.ID, userID)

	if err != nil {
		logger.Error(fmt.Sprintf("CreateTeam: failed to add creator as team member: %v", err))
		// We still return the team since it was created successfully
	}

	return &team, nil
}

func GetTeamsByUser(ctx context.Context, userID string) ([]Team, error) {
	var logger = util.NewLogger()

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("GetTeamsByUser: failed to get pool connection: %v", err))
		return nil, fmt.Errorf("something went wrong")
	}

	rows, err := db.Query(ctx,
		`SELECT t.id, t.name, t.description, t.createdBy, t.createdAt, t.updatedAt 
		FROM teams t 
		INNER JOIN team_members tm ON t.id = tm.teamID 
		WHERE tm.userID = $1 
		ORDER BY t.createdAt DESC`,
		userID)

	if err != nil {
		logger.Error(fmt.Sprintf("GetTeamsByUser: failed to query teams: %v", err))
		return nil, fmt.Errorf("something went wrong")
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.ID, &team.Name, &team.Description, &team.CreatedBy, &team.CreatedAt, &team.UpdatedAt); err != nil {
			logger.Error(fmt.Sprintf("GetTeamsByUser: failed to scan team: %v", err))
			continue
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("GetTeamsByUser: error during iteration: %v", err))
		return nil, fmt.Errorf("something went wrong")
	}

	return teams, nil
}
