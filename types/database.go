package types

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DisgoOrg/log"
	"github.com/Skye-31/WordleBot/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func SetUpDatabase(config *Config, log log.Logger) (*bun.DB, error) {
	log.Info("Setting up database")
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port)),
		pgdriver.WithUser(config.Database.User),
		pgdriver.WithPassword(config.Database.Password),
		pgdriver.WithDatabase(config.Database.DBName),
		pgdriver.WithInsecure(true),
	))
	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if config.DevMode {
		if err := db.ResetModel(context.TODO(), (*models.User)(nil)); err != nil {
			return nil, err
		}
	}
	log.Info("Database setup complete")
	return db, nil
}
