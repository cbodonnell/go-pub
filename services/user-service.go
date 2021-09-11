package services

import (
	"errors"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/repositories"
)

type service struct {
	conf config.Configuration
	repo repositories.Repository
}

func NewActivityPubService(_conf config.Configuration, _repo repositories.Repository) Service {
	return &service{
		conf: _conf,
		repo: _repo,
	}
}

func (s *service) DiscoverUserByName(name string) (models.User, error) {
	user, err := s.GetUserByName(name)
	if err != nil {
		return user, err
	}
	if !user.Discoverable {
		return user, errors.New("user is not discoverable")
	}
	// webfinger := s.resource.GenerateWebFinger(user.Name)
	return user, nil
}

func (s *service) GetUserByName(name string) (models.User, error) {
	user, err := s.repo.QueryUserByName(name)
	if err != nil {
		return user, err
	}
	// actor := s.resource.GenerateActor(user.Name)
	return user, nil
}

func (s *service) CheckUser(name string) error {
	return s.repo.CheckUser(name)
}

func (s *service) CreateUser(name string) (string, error) {
	return s.repo.CreateUser(name)
}

func (s *service) GetOutboxTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryOutboxTotalItemsByUserName(name)
}

func (s *service) GetOutboxByUserName(name string) ([]models.Activity, error) {
	return s.repo.QueryOutboxByUserName(name)
}
