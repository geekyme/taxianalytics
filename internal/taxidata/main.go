package taxidata

import (
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	influx "github.com/influxdata/influxdb1-client/v2"
)

// TaxiData refers to the message data structure received through the subscription
type TaxiData struct {
	RideID         string    `json:"ride_id"`
	PointIDX       int       `json:"point_idx"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	Timestamp      time.Time `json:"timestamp"`
	MeterReading   float64   `json:"meter_reading"`
	MeterIncrement float64   `json:"meter_increment"`
	RideStatus     string    `json:"ride_status"`
	PassengerCount int       `json:"passenger_count"`
}

// QueryResult refers to the structure we will send through a call to our handler
type QueryResult struct {
	Count int64 `json:"count"`
}

const (
	// predefined settings
	dbName      = "taxianalytics"
	seriesName  = "rides"
	bufferCount = 500
	// envs
	subNameEnv = "TAXI_SUB_NAME"
	keyEnv     = "GCLOUD_KEY"
	projectEnv = "TAXI_PROJECT"
	dbHostEnv  = "DB_HOST"
)

var (
	// TaxiSubscription points to the taxi data subscription
	TaxiSubscription *pubsub.Subscription
	// DBClient points to the influx database client
	DBClient influx.Client
)

func init() {
	subName, project, key, dbHost := os.Getenv(subNameEnv), os.Getenv(projectEnv), os.Getenv(keyEnv), os.Getenv(dbHostEnv)
	if subName == "" {
		log.Fatal(fmt.Sprintf("Must set env variable: %s", subNameEnv))
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
	TaxiSubscription, err = configureSubscription([]byte(key), project, subName)
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
