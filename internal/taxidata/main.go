package taxidata

import "time"

// TaxiData refers to the message data structure received through subscription to the taxi topic
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
	dbName      = "taxianalytics"
	seriesName  = "rides"
	bufferCount = 500
)
