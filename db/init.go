package db

import (
	"context"
	"evolve/db/connection"
	"evolve/util"
	"os"
)

func InitDb(ctx context.Context) error {
	
	logger := util.LogVar
	conn, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error("initDb: failed to get pool connection", err)
		return err
	}

	sql, err := os.ReadFile("db/scripts/init.sql")
	if err != nil {
		logger.Error("initDb: failed to read init.sql", err)
		return err
	}

	_, err = conn.Exec(ctx, string(sql))
	if err != nil {
		logger.Error("initDb: failed to execute init.sql", err)
		return err
	}

	logger.Info("initDb: database initialized")

	return nil
}
