package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// db instance
var db *pgxpool.Pool

// connect to db
func connectDb(s DataSource) *pgxpool.Pool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.Dbname)
	db, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connected to %s as %s\n", s.Dbname, s.User)
	return db
}

func queryUserByName(name string) (User, error) {
	sql := `SELECT * FROM users
	WHERE name = $1
	LIMIT 1`

	var user User
	err := db.QueryRow(context.Background(), sql, name).Scan(
		&user.ID,
		&user.Name,
		&user.Discoverable,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}
