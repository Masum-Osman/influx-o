package main

import (
	"errors"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

const (
	USERNAME string = ""
	PASSWORD string = ""
	DATABASE string = "sample_db"
)

func createPoint(params map[string]interface{}) (map[string]interface{}, error) {

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: USERNAME,
		Password: PASSWORD,
	})

	if err != nil {
		return nil, err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: DATABASE,
	})

	if err != nil {
		return nil, err
	}

	clientId := params["client_id"].(string)
	location := params["location"].(string)
	currentLoadInAmperes := params["current_load_in_amperes"]

	pt, err := client.NewPoint("power_usage", map[string]string{"client_id": clientId, "location": location},
		map[string]interface{}{"current_load_in_amperes": currentLoadInAmperes},
		time.Now())

	if err != nil {
		return nil, err
	}

	bp.AddPoint(pt)

	err = influxClient.Write(bp)

	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{"data": "Success"}

	return resp, nil

}

func getPoints() (map[string]interface{}, error) {

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: USERNAME,
		Password: PASSWORD,
	})

	if err != nil {
		return nil, err
	}

	queryString := "SELECT client_id, location, current_load_in_amperes FROM power_usage"

	q := client.NewQuery(queryString, DATABASE, "ns")

	response, err := influxClient.Query(q)

	if err != nil {
		return nil, err
	}

	err = response.Error()

	if err != nil {

		return nil, errors.New("Empty record set")

	} else {

		res := response.Results

		if len(res) == 0 {
			return nil, err
		}

		columns := response.Results[0].Series[0].Columns
		points := response.Results[0].Series[0].Values

		data := []map[string]interface{}{}

		for i := 0; i <= len(points)-1; i++ {

			record := map[string]interface{}{}

			for j := 0; j <= len(columns)-1; j++ {
				record[string(columns[j])] = points[i][j]
			}

			data = append(data, record)

		}

		resp := map[string]interface{}{"data": data}

		return resp, nil
	}
}
