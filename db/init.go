package db

import (
	"context"
	"evolve/db/connection"
	"evolve/util"
	"os"
)

func InitDb(ctx context.Context) error {
	var logger = util.NewLogger()
	conn, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error("initDb: failed to get pool connection")
		return err
	}

	sql, err := os.ReadFile("db/scripts/init.sql")
	if err != nil {
		logger.Error("initDb: failed to read init.sql")
		return err
	}

	_, err = conn.Exec(ctx, string(sql))
	if err != nil {
		logger.Error("initDb: failed to execute init.sql")
		return err
	}

	logger.Info("initDb: database initialized")

	return nil
}
