package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct{}

// TODO: Look to incorporate https://gorm.io/gorm
var (
	db *pgxpool.Pool
)

func NewPSQLRepository(source config.DataSource) Repository {
	db = connectDb(source)
	return &repository{}
}

// TODO: Make this less database dependent
func connectDb(s config.DataSource) *pgxpool.Pool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.Dbname)
	db, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	log.Printf("Connected to %s as %s\n", s.Dbname, s.User)
	return db
}

func (*repository) Close() {
	db.Close()
}

func (*repository) QueryUserByName(name string) (models.User, error) {
	sql := `SELECT * FROM users
	WHERE name = $1
	LIMIT 1`

	var user models.User
	err := db.QueryRow(context.Background(), sql, name).Scan(
		&user.ID,
		&user.Name,
		&user.Discoverable,
		&user.IRI,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (*repository) CheckUser(name string) error {
	sql := `SELECT 1 from users
	WHERE name = $1`

	var result int
	_ = db.QueryRow(context.Background(), sql, name).Scan(&result)
	if result != 1 {
		return errors.New("user does not exist")
	}
	return nil
}

func (*repository) CreateUser(name string) (string, error) {
	sql := `INSERT INTO users (name, discoverable, iri)
	VALUES ($1, true, $2)`

	iri := fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name)
	_, err := db.Exec(context.Background(), sql, name, iri)
	if err != nil {
		return iri, err
	}
	return iri, nil
}

func (*repository) QueryOutboxTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(*) FROM activities
	WHERE actor = $1`

	var count int
	err := db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
	).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (*repository) QueryOutboxByUserName(name string) ([]models.Activity, error) {
	sql := `SELECT *
	FROM activities
	WHERE actor = $1
	ORDER BY id DESC`

	rows, err := db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []models.Activity
	for rows.Next() {
		var activity_id int
		var object_id int
		activity := models.NewActivity()
		err = rows.Scan(
			&activity_id,
			&activity.Type,
			&activity.Actor,
			&object_id,
			&activity.Id,
		)
		if err != nil {
			return activities, err
		}
		object_iri, err := queryObjectIRIById(object_id)
		if err != nil {
			return activities, err
		}
		object, err := queryObjectByIRI(object_iri)
		if err != nil {
			activity.ChildObject = object_iri

		} else {
			activity.ChildObject = object
		}
		activity.To, err = queryToByActivityId(activity_id)
		if err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		return activities, err
	}
	return activities, nil
}

func queryObjectIRIById(object_id int) (string, error) {
	sql := `SELECT iri
	FROM objects WHERE id = $1;`
	var iri string
	err := db.QueryRow(context.Background(), sql, object_id).Scan(
		&iri,
	)
	if err != nil {
		return iri, err
	}
	return iri, nil
}

func queryObjectByIRI(iri string) (models.Object, error) {
	sql := `SELECT type, iri, content, attributed_to, in_reply_to
	FROM objects WHERE iri = $1;`
	object := models.NewObject()
	err := db.QueryRow(context.Background(), sql, iri).Scan(
		&object.Type,
		&object.Id,
		&object.Content,
		&object.AttributedTo,
		&object.InReplyTo,
	)
	if err != nil {
		return object, err
	}
	return object, nil
}

func queryToByActivityId(activity_id int) ([]string, error) {
	sql := `SELECT iri
	FROM activities_to
	WHERE activity_id = $1`

	rows, err := db.Query(context.Background(), sql, activity_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tos []string
	for rows.Next() {
		var to string
		err = rows.Scan(
			&to,
		)
		if err != nil {
			return tos, err
		}
		tos = append(tos, to)
	}
	err = rows.Err()
	if err != nil {
		return tos, err
	}
	return tos, nil
}
