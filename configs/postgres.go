package configs

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type (
	// DB connector interface
	DB interface {
		Connect() (db *sql.DB, err error)
	}

	// Postgres config
	Postgres struct {
		Host    string  `yaml:"host"`
		Port    int64   `yaml:"port"`
		User    string  `yaml:"user"`
		Pass    string  `yaml:"pass"`
		Dbname  string  `yaml:"dbname"`
		AppName *string `yaml:"appname"`
	}
)

var (
	DBSetMaxOpenConns    = 100
	DBSetMaxIdleConns    = 10
	DBSetConnMaxLifetime = time.Hour
)

var _ DB = &Postgres{}

// Connect raw sql connect
func (p *Postgres) Connect() (db *sql.DB, err error) {
	var appName, _ = os.Executable()
	if p.AppName != nil {
		appName = *p.AppName
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s application_name=%s sslmode=disable",
		p.Host,
		p.Port,
		p.User,
		p.Pass,
		p.Dbname,
		appName,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	db.SetMaxOpenConns(DBSetMaxOpenConns)
	db.SetMaxIdleConns(DBSetMaxIdleConns)
	db.SetConnMaxLifetime(DBSetConnMaxLifetime)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, fmt.Errorf("could not test db connection: %w", err)
	}

	return db, nil
}
