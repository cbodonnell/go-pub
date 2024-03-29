package resources

import "github.com/cheebz/go-pub/pkg/models"

type Resource interface {
	ParseResource(resource string) (string, error)
	GenerateWebFinger(name string) models.WebFinger
	GenerateActor(name string) models.Actor
	GenerateOrderedCollection(name string, endpoint string, totalItems int) models.OrderedCollection
	GenerateOrderedCollectionPage(name string, endpoint string, orderedItems []interface{}, pageNum int) models.OrderedCollectionPage
	GenerateCheckResponse(activityIRI string) models.CheckResponse
}
