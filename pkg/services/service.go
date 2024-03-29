package services

import (
	"github.com/cheebz/arb"
	"github.com/cheebz/go-pub/pkg/media"
	"github.com/cheebz/go-pub/pkg/models"
)

type Service interface {
	DiscoverUserByName(name string) (models.User, error)
	GetUserByName(name string) (models.User, error)
	CheckUser(name string) error
	CreateUser(name string) (string, error)
	GetFeedTotalItemsByUserName(name string) (int, error)
	GetFeedByUserName(name string, pageNum int) ([]models.Activity, error)
	GetInboxTotalItemsByUserName(name string) (int, error)
	GetInboxByUserName(name string, pageNum int) ([]models.Activity, error)
	GetOutboxTotalItemsByUserName(name string) (int, error)
	GetOutboxByUserName(name string, pageNum int) ([]models.Activity, error)
	GetFollowersTotalItemsByUserName(name string) (int, error)
	GetFollowersByUserName(name string, pageNum int) ([]string, error)
	GetFollowingTotalItemsByUserName(name string) (int, error)
	GetFollowingByUserName(name string, pageNum int) ([]string, error)
	GetLikedTotalItemsByUserName(name string) (int, error)
	// GetLikedByUserName(name string, pageNum int) ([]models.Object, error)
	GetLikedByUserName(name string, pageNum int) ([]string, error)
	GetActivity(ID int) (models.Activity, error)
	GetObject(ID int) (models.Object, error)
	SaveInboxActivity(activityArb arb.Arb, name string) (arb.Arb, error)
	SaveOutboxActivity(activityArb arb.Arb, name string) (arb.Arb, error)
	UploadMedia(activityArb arb.Arb, m media.Media, name string) (arb.Arb, error)
	CheckActivity(name string, activityType string, objectIRI string) string
}
