package taxidata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/fatih/structs"
	influx "github.com/influxdata/influxdb1-client/v2"
	"google.golang.org/api/option"
)

const (
	topicEnv   = "TAXI_TOPIC"
	keyEnv     = "GCLOUD_KEY"
	projectEnv = "TAXI_PROJECT"
	dbHostEnv  = "DB_HOST"
)

var (
	// TaxiSubscription points to the active subscription we have on the taxi topic
	TaxiSubscription *pubsub.Subscription
	// DBClient points to the influx database client
	DBClient influx.Client
)

func init() {
	topic, project, key, dbHost := os.Getenv(topicEnv), os.Getenv(projectEnv), os.Getenv(keyEnv), os.Getenv(dbHostEnv)
	if topic == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", topicEnv))
	}
	if project == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", projectEnv))
	}
	if key == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", keyEnv))
	}
	if dbHost == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", dbHostEnv))
	}

	var err error
	TaxiSubscription, err = configureSubscription([]byte(key), project, topic)
	if err != nil {
		log.Fatal(err)
	}

	DBClient, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr: dbHost,
	})

	if err != nil {
		log.Fatal(err)
	}
}

// Subscribe takes in a subscription and runs callback functions
func Subscribe(subscription *pubsub.Subscription) {
	ctx := context.Background()
	bufCh := make(chan TaxiData, bufferCount)
	go batchWriter(bufCh)

	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var data TaxiData
		var err error
		if err = json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Nack()
			return
		}

		bufCh <- data
		msg.Ack()
	})
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

	sub := client.Subscription(topic)
	// This is needed otherwise we risk this subscriber not Ack-ing a message
	// We should scale horizontally to allow more workers to consume the message
	sub.ReceiveSettings.MaxOutstandingMessages = bufferCount

	return sub, nil
}

func batchWriter(ch chan TaxiData) {
	for {
		var items []TaxiData

		items = append(items, <-ch)

		// we batch write at half the channel buffer length each time
		// so that the channel can continue accumulating items during the write
		remains := (bufferCount / 2) - 1

		for i := 0; i < remains; i++ {
			select {
			case item := <-ch:
				items = append(items, item)
			default:
				break
			}
		}

		if err := insertPoints(DBClient, dbName, items); err != nil {
			log.Printf("Error inserting data points in batch: %v", err)
		}
	}
}

func insertPoints(c influx.Client, dbName string, items []TaxiData) error {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: dbName,
	})
	if err != nil {
		return err
	}

	for _, data := range items {
		// Create a point and add to batch
		tags := map[string]string{}
		fields := structs.Map(data)

		pt, err := influx.NewPoint(seriesName, tags, fields, data.Timestamp)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}
