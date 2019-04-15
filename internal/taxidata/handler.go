package taxidata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Handler writes current state of the taxi data to the requester
func Handler(w http.ResponseWriter, r *http.Request) {
	res, err := getAllInLast(DBClient, "1h")

	if err != nil {
		log.Printf("Error fetching ride statistics: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func getAllInLast(c influx.Client, timeframe string) (*QueryResult, error) {
	// the dataset shows that a ride status is either 'enroute' or 'pickup'
	// there is usually only 1 pickup per ride,
	// so we approximate trips per X using the following expression
	expression := fmt.Sprintf("select count(distinct(RideID)) from %s where RideStatus = 'pickup' and time > now() - %s", seriesName, timeframe)
	q := influx.NewQuery(expression, dbName, "ns")
	response, err := c.Query(q)

	if err != nil {
		return nil, err
	}

	if err = response.Error(); err != nil {
		return nil, err
	}

	if response.Results[0].Series == nil {
		return &QueryResult{Count: 0}, nil
	} else {
		num, err := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		return &QueryResult{Count: num}, nil
	}
}
