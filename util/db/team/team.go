package team

import (
	"context"
	"evolve/util"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TeamExistsByNameForUser checks if a team with the given name exists for a specific user (createdBy)
func TeamExistsByNameForUser(ctx context.Context, teamName string, createdBy string, db *pgxpool.Pool) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM team WHERE teamName = $1 AND createdBy = $2)", teamName, createdBy).Scan(&exists)
	if err != nil {
		util.SharedLogger.Error(fmt.Sprintf("TeamExistsByNameForUser: failed to check team existence: %v", err), err)
		return false, fmt.Errorf("something went wrong")
	}
	return exists, nil
}

// TeamExistsByName checks if a team with the given name exists (any user)
func TeamExistsByName(ctx context.Context, teamName string, db *pgxpool.Pool) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM team WHERE teamName = $1)", teamName).Scan(&exists)
	if err != nil {
		util.SharedLogger.Error(fmt.Sprintf("TeamExistsByName: failed to check team existence: %v", err), err)
		return false, fmt.Errorf("something went wrong")
	}
	return exists, nil
}

// IsUserTeamMember checks if a user is a member of a team (either as creator or as a member in teamMembers table)
func IsUserTeamMember(ctx context.Context, teamName string, userID string, db *pgxpool.Pool) (bool, error) {
	var isMember bool
	err := db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM team T WHERE T.teamName = $1 AND (T.createdBy = $2 OR EXISTS(SELECT 1 FROM teamMembers TM WHERE TM.teamID = T.teamID AND TM.memberId = $3::uuid)))",
		teamName, userID, userID).Scan(&isMember)
	if err != nil {
		util.SharedLogger.Error(fmt.Sprintf("IsUserTeamMember: failed to check team membership: %v", err), err)
		return false, fmt.Errorf("something went wrong")
	}
	return isMember, nil
}
