package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/cheebz/arb"
	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/go-pub/workers"
)

type ActivityPubService struct {
	conf   config.Configuration
	repo   repositories.Repository
	worker workers.Worker
}

func NewActivityPubService(_conf config.Configuration, _repo repositories.Repository, _worker workers.Worker) Service {
	return &ActivityPubService{
		conf:   _conf,
		repo:   _repo,
		worker: _worker,
	}
}

func (s *ActivityPubService) DiscoverUserByName(name string) (models.User, error) {
	user, err := s.GetUserByName(name)
	if err != nil {
		return user, err
	}
	if !user.Discoverable {
		return user, errors.New("user is not discoverable")
	}
	return user, nil
}

func (s *ActivityPubService) GetUserByName(name string) (models.User, error) {
	return s.repo.QueryUserByName(name)
}

func (s *ActivityPubService) CheckUser(name string) error {
	return s.repo.CheckUser(name)
}

func (s *ActivityPubService) CreateUser(name string) (string, error) {
	return s.repo.CreateUser(name)
}

func (s *ActivityPubService) GetInboxTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryInboxTotalItemsByUserName(name)
}

func (s *ActivityPubService) GetInboxByUserName(name string, pageNum int) ([]models.Activity, error) {
	return s.repo.QueryInboxByUserName(name, pageNum)
}

func (s *ActivityPubService) GetOutboxTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryOutboxTotalItemsByUserName(name)
}

func (s *ActivityPubService) GetOutboxByUserName(name string, pageNum int) ([]models.Activity, error) {
	return s.repo.QueryOutboxByUserName(name, pageNum)
}

func (s *ActivityPubService) GetFollowersTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryFollowersTotalItemsByUserName(name)
}

func (s *ActivityPubService) GetFollowersByUserName(name string, pageNum int) ([]string, error) {
	return s.repo.QueryFollowersByUserName(name, pageNum)
}

func (s *ActivityPubService) GetFollowingTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryFollowingTotalItemsByUserName(name)
}

func (s *ActivityPubService) GetFollowingByUserName(name string, pageNum int) ([]string, error) {
	return s.repo.QueryFollowingByUserName(name, pageNum)
}

func (s *ActivityPubService) GetLikedTotalItemsByUserName(name string) (int, error) {
	return s.repo.QueryLikedTotalItemsByUserName(name)
}

// func (s *ActivityPubService) GetLikedByUserName(name string, pageNum int) ([]models.Object, error) {
func (s *ActivityPubService) GetLikedByUserName(name string, pageNum int) ([]string, error) {
	return s.repo.QueryLikedByUserName(name, pageNum)
}

func (s *ActivityPubService) GetActivity(ID int) (models.Activity, error) {
	return s.repo.QueryActivity(ID)
}

func (s *ActivityPubService) GetObject(ID int) (models.Object, error) {
	return s.repo.QueryObject(ID)
}

func (s *ActivityPubService) SaveInboxActivity(activityArb arb.Arb, name string) (arb.Arb, error) {
	activityIRI, err := activitypub.GetIRI(activityArb)
	if err != nil {
		return activityArb, err
	}
	actorArb, err := activitypub.FindProp(activityArb, "actor", activitypub.AcceptHeaders)
	if err != nil {
		return activityArb, err
	}
	actorIRI, err := activitypub.GetIRI(actorArb)
	if err != nil {
		return activityArb, err
	}
	objectArb, err := activitypub.FindProp(activityArb, "object", activitypub.AcceptHeaders)
	if err != nil {
		return activityArb, err
	}
	objectIRI, err := activitypub.GetIRI(objectArb)
	if err != nil {
		return activityArb, err
	}
	activityType, err := activitypub.GetType(activityArb)
	if err != nil {
		return activityArb, err
	}
	recipient := fmt.Sprintf("%s://%s/%s/%s", s.conf.Protocol, s.conf.ServerName, s.conf.Endpoints.Users, name)
	switch activityType {
	// case "Create":
	// 	_, err = s.repo.CreateInboxActivity(activityArb, objectArb, actorIRI.String(), name)
	// 	if err != nil {
	// 		return activityArb, err
	// 	}
	case "Create", "Announce", "Like", "Undo", "Accept":
		_, err = s.repo.CreateInboxReferenceActivity(activityArb, objectIRI.String(), actorIRI.String(), name)
		if err != nil {
			return activityArb, err
		}
	case "Follow":
		if objectIRI.String() != recipient {
			return activityArb, errors.New("wrong inbox")
		}
		_, err = s.repo.CreateInboxReferenceActivity(activityArb, recipient, actorIRI.String(), name)
		if err != nil {
			return activityArb, err
		}
		responseArb, err := activitypub.NewActivityArbReference(activityIRI.String(), "Accept")
		if err != nil {
			return activityArb, err
		}
		responseArb["actor"] = recipient
		responseArb, err = s.repo.CreateOutboxReferenceActivity(responseArb, name)
		if err != nil {
			return activityArb, err
		}
		s.worker.GetChannel() <- models.Federation{Name: name, Recipient: actorIRI.String(), Activity: responseArb}
	case "Delete":
		// TODO: DeleteActivity
		attributedTo, err := objectArb.GetString("attributedTo")
		if err != nil {
			return activityArb, err
		}
		if actorIRI.String() != attributedTo {
			return activityArb, errors.New("not your object")
		}
		// Add Delete Activity to collection
		// Replace activity with Tombstone (or delete all together?)
	default:
		return activityArb, errors.New("unsupported activity type")
	}
	return activityArb, nil
}

func (s *ActivityPubService) SaveOutboxActivity(activityArb arb.Arb, name string) (arb.Arb, error) {
	objectArb, err := activitypub.FindProp(activityArb, "object", activitypub.AcceptHeaders)
	if err != nil {
		return activityArb, err
	}
	actor := fmt.Sprintf("%s://%s/%s/%s", s.conf.Protocol, s.conf.ServerName, s.conf.Endpoints.Users, name)
	activityArb["actor"] = actor
	activityType, err := activitypub.GetType(activityArb)
	if err != nil {
		return activityArb, err
	}
	switch activityType {
	case "Create":
		objectArb["attributedTo"] = actor
		activityArb, err = s.repo.CreateOutboxActivity(activityArb, objectArb, name)
		if err != nil {
			return activityArb, err
		}
	case "Follow":
		activityArb, err = s.repo.CreateOutboxReferenceActivity(activityArb, name)
		if err != nil {
			return activityArb, err
		}
		objectIRI, err := activitypub.GetIRI(objectArb)
		if err != nil {
			return activityArb, err
		}
		// check if the recipient is internal
		if objectIRI.Host == s.conf.ServerName {
			// if so, generate and federate an accept
			activityIRI, err := activitypub.GetIRI(activityArb)
			if err != nil {
				return activityArb, err
			}
			responseArb, err := activitypub.NewActivityArbReference(activityIRI.String(), "Accept")
			if err != nil {
				return activityArb, err
			}
			responseArb["actor"] = objectIRI.String()
			responseArb, err = s.repo.CreateOutboxReferenceActivity(responseArb, name)
			if err != nil {
				return activityArb, err
			}
			s.worker.GetChannel() <- models.Federation{Name: name, Recipient: actor, Activity: responseArb}
		}
	case "Like", "Undo":
		activityArb, err = s.repo.CreateOutboxReferenceActivity(activityArb, name)
		if err != nil {
			return activityArb, err
		}
	case "Delete":
		// TODO: DeleteActivity
		attributedTo, err := objectArb.GetString("attributedTo")
		if err != nil {
			return activityArb, err
		}
		if actor != attributedTo {
			return activityArb, errors.New("not your object")
		}
		activityArb, err = s.repo.DeleteActivity(activityArb, name)
		if err != nil {
			return activityArb, err
		}
	default:
		return activityArb, errors.New("unsupported activity type")
	}
	// Get recipients
	recipients, err := activitypub.GetRecipients(activityArb, "to")
	if err != nil {
		log.Println(err)
	}
	// Deliver to recipients
	for _, recipient := range recipients {
		s.worker.GetChannel() <- models.Federation{Name: name, Recipient: recipient.String(), Activity: activityArb}
	}
	return activityArb, nil
}
