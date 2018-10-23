package database

import (
	"fmt"
	"time"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/jinzhu/gorm"

	// Drivers
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/lib/pq"
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// NewDatabase returns a gorm.DB struct, gorm.DB.DB() returns a database handle
// see http://golang.org/pkg/database/sql/#DB
func NewDatabase(cnf *config.Config) (*gorm.DB, error) {
	// Connection args
	// see https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	var args string
	switch cnf.Database.Type {
	case "postgres":
		args = fmt.Sprintf(
			"sslmode=disable host=%s port=%d user=%s password='%s' dbname=%s",
			cnf.Database.Host,
			cnf.Database.Port,
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.DatabaseName,
		)
	case "mssql":
		args = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.Host,
			cnf.Database.Port,
			cnf.Database.DatabaseName,
		)
	}

	if args == "" {
		// Database type not supported
		return nil, fmt.Errorf("Database type %s not suppported", cnf.Database.Type)
	}

	db, err := gorm.Open(cnf.Database.Type, args)
	if err != nil {
		return db, err
	}

	// Max idle connections
	db.DB().SetMaxIdleConns(cnf.Database.MaxIdleConns)

	// Max open connections
	db.DB().SetMaxOpenConns(cnf.Database.MaxOpenConns)

	// Database logging
	db.LogMode(cnf.IsDevelopment)

	return db, nil
}
