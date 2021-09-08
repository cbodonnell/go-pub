package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/cheebz/go-pub/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct{}

var (
	db *pgxpool.Pool
)

func NewPSQLRepository(source *pgxpool.Pool) Repository {
	db = source
	return &repository{}
}

// TODO: Move this to a database package with Connect and Close methods
func ConnectDb(s models.DataSource) *pgxpool.Pool {
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
