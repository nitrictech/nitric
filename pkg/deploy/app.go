package deploy

import (
	"fmt"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// Collections and stores information about a nitric application
// And it's dependencies
type App struct {
	apis          map[string][]*pb.ApiWorker
	subscriptions map[string]*pb.SubscriptionWorker
	schedules     map[string]*pb.ScheduleWorker
	buckets       map[string]*pb.BucketResource
	topics        map[string]*pb.TopicResource
	collections   map[string]*pb.CollectionResource
	queues        map[string]*pb.QueueResource
	policies      []*pb.PolicyResource
}

func (a *App) String() string {
	return fmt.Sprintf(`
	  apis: %+v,
	  subscriptions: %+v,
		schedules: %+v,
		buckets: %+v,
		topics: %+v,
		collections: %+v,
		queues: %+v,
		policies: %+v,
	`, a.apis, a.subscriptions, a.schedules, a.buckets, a.topics, a.queues, a.collections, a.policies)
}

func (a *App) AddPolicy(p *pb.PolicyResource) {
	a.policies = append(a.policies, p)
}

func (a *App) AddApiHandler(aw *pb.ApiWorker) {
	if a.apis[aw.Api] != nil {
		a.apis[aw.Api] = make([]*pb.ApiWorker, 0)
	}

	a.apis[aw.Api] = append(a.apis[aw.Api], aw)
}

func (a *App) AddSubscriptionHandler(sw *pb.SubscriptionWorker) error {
	// TODO: Determine if this subscription handler has a write policy to the same topic
	if a.subscriptions[sw.Topic] != nil {
		// return a new error
		return fmt.Errorf("subscription already declared for topic %s, only one subscription per topic is allowed per application", sw.Topic)
	}

	// This maps to a trigger worker for this application
	a.subscriptions[sw.Topic] = sw

	return nil
}

func (a *App) AddScheduleHandler(sw *pb.ScheduleWorker) error {
	if a.schedules[sw.Key] != nil {
		return fmt.Errorf("schedule %s already exists", sw.Key)
	}

	a.schedules[sw.Key] = sw

	return nil
}

func (a *App) AddBucket(name string, b *pb.BucketResource) {
	a.buckets[name] = b
}

func (a *App) AddTopic(name string, t *pb.TopicResource) {
	a.topics[name] = t
}

func (a *App) AddCollection(name string, c *pb.CollectionResource) {
	a.collections[name] = c
}

func (a *App) AddQueue(name string, q *pb.QueueResource) {
	a.queues[name] = q
}

func NewApp() *App {
	return &App{
		apis:          make(map[string][]*pb.ApiWorker),
		subscriptions: make(map[string]*pb.SubscriptionWorker),
		schedules:     make(map[string]*pb.ScheduleWorker),
		buckets:       make(map[string]*pb.BucketResource),
		topics:        make(map[string]*pb.TopicResource),
		collections:   make(map[string]*pb.CollectionResource),
		queues:        make(map[string]*pb.QueueResource),
		policies:      make([]*pb.PolicyResource, 0),
	}
}
