package repositories

import "github.com/cheebz/go-pub/models"

type Repository interface {
	Close()
	QueryUserByName(name string) (models.User, error)
	CheckUser(name string) error
	CreateUser(name string) (string, error)
}
