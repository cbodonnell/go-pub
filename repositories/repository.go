package repositories

import "github.com/cheebz/go-pub/models"

type Repository interface {
	QueryUserByName(name string) (models.User, error)
}
