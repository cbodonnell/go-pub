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
