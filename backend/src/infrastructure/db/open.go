// infrastructure/db/open.go
package db

import (
	"database/sql"

	"word_app/backend/ent"

	sqldrv "entgo.io/ent/dialect/sql"
)

func OpenShared(dsn string, driverName string) (*sql.DB, *ent.Client, error) {
	// 例: driverName = "pgx" / "postgres"
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, err
	}
	drv := sqldrv.OpenDB(driverName, db) // または sqldrv.NewDriver(db)
	client := ent.NewClient(ent.Driver(drv))
	return db, client, nil
}
