package repositories

import (
	"github.com/cheebz/arb"
	"github.com/cheebz/go-pub/models"
)

type Repository interface {
	Close()
	QueryUserByName(name string) (models.User, error)
	CheckUser(name string) error
	CreateUser(name string) (string, error)
	QueryInboxTotalItemsByUserName(name string) (int, error)
	QueryInboxByUserName(name string) ([]models.Activity, error)
	QueryOutboxTotalItemsByUserName(name string) (int, error)
	QueryOutboxByUserName(name string) ([]models.Activity, error)
	QueryFollowersTotalItemsByUserName(name string) (int, error)
	QueryFollowersByUserName(name string) ([]string, error)
	QueryFollowingTotalItemsByUserName(name string) (int, error)
	QueryFollowingByUserName(name string) ([]string, error)
	QueryLikedTotalItemsByUserName(name string) (int, error)
	QueryLikedByUserName(name string) ([]models.Object, error)
	QueryActivity(ID int) (models.Activity, error)
	QueryObject(ID int) (models.Object, error)
	CreateInboxActivity(activityArb arb.Arb, objectArb arb.Arb, actor string, recipient string) (arb.Arb, error)
	CreateInboxReferenceActivity(activityArb arb.Arb, object string, actor string, recipient string) (arb.Arb, error)
	CreateOutboxActivity(activityArb arb.Arb, objectArb arb.Arb) (arb.Arb, error)
	CreateOutboxReferenceActivity(activityArb arb.Arb) (arb.Arb, error)
	ActivityToExists(activityIRI string, recipientIRI string) bool
	AddActivityTo(activityIRI string, recipient string) error
}
