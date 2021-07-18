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
		&user.URL,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func queryInboxTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(ps.*)
	FROM notes as ps
	INNER JOIN activities as act
	ON act.object_id = ps.id
	INNER JOIN activities_to as act_to
	ON act_to.activity_id = act.id
	INNER JOIN users as us
	ON us.id = act_to.to
	WHERE us.name = $1`

	var count int
	err := db.QueryRow(context.Background(), sql, name).Scan(
		&count,
	)
	if err != nil {
		return count, err
	}
	return count, nil
}

func queryInboxByUserName(name string) ([]Note, error) {
	sql := `SELECT ps.*, act.id, act.user_name, act.type, us.url
	FROM notes as ps
	INNER JOIN activities as act
	ON act.object_id = ps.id
	INNER JOIN activities_to as act_to
	ON act_to.activity_id = act.id
	INNER JOIN users as us
	ON us.id = act_to.to
	WHERE us.name = $1
	ORDER BY act.id DESC`

	rows, err := db.Query(context.Background(), sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []Note
	for rows.Next() {
		var note Note
		var activity Activity
		var to string
		err = rows.Scan(
			&note.ID,
			&note.UserName,
			&note.Content,
			&activity.ID,
			&activity.UserName,
			&activity.Type,
			&to,
		)
		if err != nil {
			return notes, err
		}
		// TODO: Do this a better way... maybe a second query?
		activity.To = []string{to}
		note.Activity = activity
		notes = append(notes, note)
	}
	err = rows.Err()
	if err != nil {
		return notes, err
	}
	return notes, nil
}

func queryOutboxTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(*) FROM notes
	WHERE user_name = $1`

	var count int
	err := db.QueryRow(context.Background(), sql, name).Scan(
		&count,
	)
	if err != nil {
		return count, err
	}
	return count, nil
}

func queryOutboxByUserName(name string) ([]Note, error) {
	sql := `SELECT ps.*, act.id, act.user_name, act.type
	FROM notes as ps
	INNER JOIN activities as act
	ON act.object_id = ps.id
	WHERE ps.user_name = $1
	ORDER BY act.id DESC`

	rows, err := db.Query(context.Background(), sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []Note
	for rows.Next() {
		var note Note
		var activity Activity
		err = rows.Scan(
			&note.ID,
			&note.UserName,
			&note.Content,
			&activity.ID,
			&activity.UserName,
			&activity.Type,
		)
		if err != nil {
			return notes, err
		}
		note.Activity = activity
		notes = append(notes, note)
	}
	err = rows.Err()
	if err != nil {
		return notes, err
	}
	return notes, nil
}
