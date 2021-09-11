package services

import "github.com/cheebz/go-pub/models"

type Service interface {
	DiscoverUserByName(name string) (models.User, error)
	GetUserByName(name string) (models.User, error)
	CheckUser(name string) error
	CreateUser(name string) (string, error)
	GetOutboxTotalItemsByUserName(name string) (int, error)
	GetOutboxByUserName(name string) ([]models.Activity, error)
}
