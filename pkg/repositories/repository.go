package repositories

import (
	"github.com/cheebz/arb"
	"github.com/cheebz/go-pub/pkg/models"
)

type Repository interface {
	Close()
	QueryUserByName(name string) (models.User, error)
	CheckUser(name string) error
	CreateUser(name string) (string, error)
	QueryFeedTotalItemsByUserName(name string) (int, error)
	QueryFeedByUserName(name string, pageNum int) ([]models.Activity, error)
	QueryInboxTotalItemsByUserName(name string) (int, error)
	QueryInboxByUserName(name string, pageNum int) ([]models.Activity, error)
	QueryOutboxTotalItemsByUserName(name string) (int, error)
	QueryOutboxByUserName(name string, pageNum int) ([]models.Activity, error)
	QueryFollowersTotalItemsByUserName(name string) (int, error)
	QueryFollowersByUserName(name string, pageNum int) ([]string, error)
	QueryFollowingTotalItemsByUserName(name string) (int, error)
	QueryFollowingByUserName(name string, pageNum int) ([]string, error)
	QueryLikedTotalItemsByUserName(name string) (int, error)
	QueryLikedByUserName(name string, pageNum int) ([]string, error)
	QueryActivity(ID int) (models.Activity, error)
	QueryObject(ID int) (models.Object, error)
	CreateInboxActivity(activityArb arb.Arb, objectArb arb.Arb, actor string, name string) (arb.Arb, error)
	CreateInboxReferenceActivity(activityArb arb.Arb, object string, actor string, name string) (arb.Arb, error)
	CreateOutboxActivity(activityArb arb.Arb, objectArb arb.Arb, name string) (arb.Arb, error)
	CreateOutboxReferenceActivity(activityArb arb.Arb, name string) (arb.Arb, error)
	ActivityToExists(activityIRI string, recipientIRI string) bool
	AddActivityTo(activityIRI string, recipient string) error
	DeleteActivity(activityArb arb.Arb, name string) (arb.Arb, error)
	GetObjectFilesByIRI(objectIRI string) ([]string, error)
	PurgeUnusedFiles() error
	CheckActivity(name string, activityType string, objectIRI string) string
}
