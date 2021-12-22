package deploy

import (
	"fmt"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// Collections and stores information about a nitric application
// And it's dependencies
type App struct {
	apis          map[string]*ApiBuilder
	subscriptions map[string]*Subscription
	schedules     map[string]*Schedule
	buckets       map[string]*Bucket
	topics        map[string]*Topic
	collections   map[string]*Collection
	queues        map[string]*Queue
}

func (a *App) AddApiHandler(rw *pb.RouteWorker) {
	if a.apis[rw.Api] != nil {
		a.apis[rw.Api] = NewApiBuilder()
	}

	a.apis[rw.Api].AddRouteHandler(rw.Api, &RouteHandler{
		path: rw.Path,
		methods: rw.Methods,
	})

}

func (a *App) AddSubscriptionHandler(sw *pb.SubscriptionWorker) error {
	// See if there is a topic publish permission for this application already
	if a.topics[sw.Topic] != nil {

	}

	if a.subscriptions[sw.Topic] != nil {
		// return a new error
		return fmt.Errorf("subscription already declared for topic %s, only one subscription per topic is allowed per application", sw.Topic)
	}

	// This maps to a trigger worker for this application
	a.subscriptions[sw.Topic] = &Subscription{}

	return nil
}

func (a *App) AddScheduleHandler(sw *pb.ScheduleWorker) error {
	if a.schedules[sw.Key] != nil {
		return fmt.Errorf("schedule %s already exists", sw.Key)
	}

	a.schedules[sw.Key] = &Schedule{}

	return nil
}

func (a *App) AddBucket(b *pb.BucketResource) {
	if a.buckets[b.Name] != nil {

	}

	// TODO: Handle duplicate resource declarations (either by merge or failure)
}

func (a *App) AddTopic(t *pb.TopicResource) {
	if a.topics[t.Name] != nil {

	}

	// TODO: Handle duplicate resource declarations (either by merge or failure)
}

func (a *App) AddCollection(c *pb.CollectionResource) {
	if a.collections[c.Name] != nil {

	}
}

func (a *App) AddQueue(q *pb.QueueResource) {
	queue := a.queues[q.Name]

	if queue == nil {
		// Handle error or merging
		queue = &Queue{}
	}

	// Add permissions

	// Update queue reference
	a.queues[c.Name] = queue

func NewApp() *App {
	return &App{
		apis:          make(map[string]*ApiBuilder),
		subscriptions: make(map[string]*Subscription),
		schedules:     make(map[string]*Schedule),
		buckets:       make(map[string]*Bucket),
		topics:        make(map[string]*Topic),
		collections:   make(map[string]*Collection),
		queues:        make(map[string]*Queue),
	}
}
