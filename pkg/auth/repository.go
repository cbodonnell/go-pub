package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cheebz/go-pub/pkg/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

// User is the auth-layer view of a user row (password included).
type User struct {
	ID       int
	Username string
	Password string
	Created  time.Time
	UUID     string
}

// Repository defines the DB operations needed by the auth service.
type Repository interface {
	GetUserByID(id int) (User, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(username, passwordHash, userUUID string) (User, error)
	UpdatePassword(userID int, passwordHash string) error
	SaveRefresh(userID int, jti string, maxAgeSecs int) error
	ValidateRefresh(userID int, jti string) error
	InvalidateRefresh(jti string) error
	DeleteAllRefresh(userID int) error
	DeleteExpiredRefresh() error
}

// PSQLRepository is the PostgreSQL implementation.
type PSQLRepository struct {
	db   *pgxpool.Pool
	conf config.Configuration
}

func NewPSQLRepository(conf config.Configuration) Repository {
	return &PSQLRepository{
		db:   connectDB(conf.Db),
		conf: conf,
	}
}

func connectDB(ds config.DataSource) *pgxpool.Pool {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		ds.Host, ds.Port, ds.User, ds.Password, ds.Dbname)
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("auth: unable to connect to database: %v", err)
	}
	return pool
}

func (r *PSQLRepository) GetUserByID(id int) (User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, created, uuid FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Created, &u.UUID)
	return u, err
}

func (r *PSQLRepository) GetUserByUsername(username string) (User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, created, uuid FROM users WHERE username = $1`, username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Created, &u.UUID)
	return u, err
}

func (r *PSQLRepository) CreateUser(username, passwordHash, userUUID string) (User, error) {
	var u User
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password, created, uuid, discoverable, iri)
		 VALUES ($1, $2, NOW(), $3, true, '')
		 RETURNING id, username, password, created, uuid`,
		username, passwordHash, userUUID,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Created, &u.UUID)
	return u, err
}

func (r *PSQLRepository) UpdatePassword(userID int, passwordHash string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET password = $1 WHERE id = $2`, passwordHash, userID)
	return err
}

func (r *PSQLRepository) SaveRefresh(userID int, jti string, maxAgeSecs int) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO user_refresh (user_id, jti, expires)
		 VALUES ($1, $2, current_timestamp + ($3 || ' seconds')::interval)`,
		userID, jti, fmt.Sprintf("%d", maxAgeSecs))
	return err
}

func (r *PSQLRepository) ValidateRefresh(userID int, jti string) error {
	var found int
	err := r.db.QueryRow(context.Background(),
		`SELECT 1 FROM user_refresh WHERE user_id = $1 AND jti = $2 AND expires > current_timestamp`,
		userID, jti,
	).Scan(&found)
	if err != nil || found != 1 {
		if err == nil {
			err = fmt.Errorf("invalid refresh token")
		}
		return err
	}
	return nil
}

func (r *PSQLRepository) InvalidateRefresh(jti string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE user_refresh SET expires = current_timestamp WHERE jti = $1`, jti)
	return err
}

func (r *PSQLRepository) DeleteAllRefresh(userID int) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM user_refresh WHERE user_id = $1`, userID)
	return err
}

func (r *PSQLRepository) DeleteExpiredRefresh() error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM user_refresh WHERE expires < current_timestamp`)
	return err
}
