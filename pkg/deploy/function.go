package deploy

import (
	"fmt"
	"sync"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// Function - Stores information about a Nitric Function, and it's dependencies
type Function struct {
	apis          map[string][]*pb.ApiWorker
	subscriptions map[string]*pb.SubscriptionWorker
	schedules     map[string]*pb.ScheduleWorker
	buckets       map[string]*pb.BucketResource
	topics        map[string]*pb.TopicResource
	collections   map[string]*pb.CollectionResource
	queues        map[string]*pb.QueueResource
	policies      []*pb.PolicyResource
	lock		sync.RWMutex
}

func (a *Function) String() string {
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

// AddPolicy - Adds an access policy dependency to the function
func (a *Function) AddPolicy(p *pb.PolicyResource) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.policies = append(a.policies, p)
}

func (a *Function) AddApiHandler(aw *pb.ApiWorker) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.apis[aw.Api] != nil {
		a.apis[aw.Api] = make([]*pb.ApiWorker, 0)
	}

	a.apis[aw.Api] = append(a.apis[aw.Api], aw)
}

// AddSubscriptionHandler - registers a handler in the function that subscribes to a topic of events
func (a *Function) AddSubscriptionHandler(sw *pb.SubscriptionWorker) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	// TODO: Determine if this subscription handler has a write policy to the same topic
	if a.subscriptions[sw.Topic] != nil {
		// return a new error
		return fmt.Errorf("subscription already declared for topic %s, only one subscription per topic is allowed per application", sw.Topic)
	}

	// This maps to a trigger worker for this application
	a.subscriptions[sw.Topic] = sw

	return nil
}

// AddScheduleHandler - registers a handler in the function that runs on a schedule
func (a *Function) AddScheduleHandler(sw *pb.ScheduleWorker) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.schedules[sw.Key] != nil {
		return fmt.Errorf("schedule %s already exists", sw.Key)
	}

	a.schedules[sw.GetKey()] = sw

	return nil
}

// AddBucket - adds a storage bucket dependency to the function
func (a *Function) AddBucket(name string, b *pb.BucketResource) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.buckets[name] = b
}

// AddTopic - adds a pub/sub topic dependency to the function
func (a *Function) AddTopic(name string, t *pb.TopicResource) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.topics[name] = t
}

// AddCollection - adds a document database collection dependency to the function
func (a *Function) AddCollection(name string, c *pb.CollectionResource) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.collections[name] = c
}

// AddQueue - adds a queue dependency to the function
func (a *Function) AddQueue(name string, q *pb.QueueResource) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.queues[name] = q
}

// NewFunction - creates a new Nitric Function, ready to register handlers and dependencies.
func NewFunction() *Function {
	return &Function{
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
