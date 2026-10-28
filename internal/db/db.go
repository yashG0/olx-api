package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(databaseUrl string) (*sql.DB, error) {
	// create the database Pool instance using sql.Open we have passed the sql driver and DB url
	db, err := sql.Open("pgx", databaseUrl)

	// if db failed to connect raise the error
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	// configure database credentials
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	// if db failed to connect within Time exit the program
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// checking the DB ping
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %v", err)
	}
	return db, nil
}
