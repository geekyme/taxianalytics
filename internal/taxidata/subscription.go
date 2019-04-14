package taxidata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	// TaxiSubscription points to the active subscription we have on the taxi topic
	TaxiSubscription *pubsub.Subscription
)

const topicEnv = "TAXI_TOPIC"
const projectEnv = "TAXI_PROJECT"
const keyEnv = "GCLOUD_KEY"

var mutex = &sync.Mutex{}
var state = []TaxiData{}

func init() {
	topic, project, key := os.Getenv(topicEnv), os.Getenv(projectEnv), os.Getenv(keyEnv)
	if topic == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", topicEnv))
	}
	if project == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", projectEnv))
	}
	if key == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", keyEnv))
	}

	var err error
	TaxiSubscription, err = configureSubscription([]byte(key), project, topic)
	if err != nil {
		log.Fatal(err)
	}
}

func configureSubscription(key []byte, project, topic string) (*pubsub.Subscription, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project, option.WithCredentialsJSON(key))
	if err != nil {
		return nil, err
	}

	if _, err := client.Topic(topic).Exists(ctx); err != nil {
		return nil, err
	}

	subscription := client.Subscription(topic)

	return subscription, nil
}

// Subscribe takes in a subscription and runs callback functions
func Subscribe(subscription *pubsub.Subscription) {
	ctx := context.Background()
	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var data TaxiData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Ack()
			return
		}

		mutex.Lock()
		state = append(state, data)
		mutex.Unlock()

		msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Handler writes current state of the taxi data to the requester
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Count: %d\n", len(state))
	json.NewEncoder(w).Encode(state)
}
