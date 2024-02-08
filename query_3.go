package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Define the InfluxDB connection parameters
	influxURL := "http://localhost:8086"
	influxToken := "your_token"
	influxOrg := "your_org"
	influxBucket := "your_bucket_name"

	// Create a new InfluxDB client
	client := influxdb2.NewClient(influxURL, influxToken)

	// Instantiate a new Flux query service
	queryAPI := client.QueryAPI(influxOrg)

	// Define the start and end time for the query
	startTime := time.Now().Add(-30 * 24 * time.Hour)
	endTime := time.Now()

	// Construct Flux queries for each measurement
	queries := []string{
		fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "classes")
		`, influxBucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)),
		fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "class_participants")
		`, influxBucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)),
		fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "quiz_participants")
		`, influxBucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)),
		fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] == "exam_participants")
		`, influxBucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)),
	}

	// Execute Flux queries and print results
	for _, query := range queries {
		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return
		}
		defer result.Close()

		for result.Next() {
			// Process and print query results as needed
			// Note: Modify this section based on your requirements
			fmt.Println(result.Record())
		}
		if err := result.Err(); err != nil {
			fmt.Println("Error iterating query results:", err)
			return
		}
	}

	// Close the InfluxDB client
	client.Close()
}
