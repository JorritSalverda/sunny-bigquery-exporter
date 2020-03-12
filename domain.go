package main

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type BigQueryMeasurement struct {
	Location   string         `bigquery:"location"`
	MeasuredAt time.Time      `bigquery:"measured_at"`
	Zones      []BigQueryZone `bigquery:"zones"`
	InsertedAt time.Time      `bigquery:"inserted_at"`
}

type BigQueryZone struct {
	Zone              string               `bigquery:"location"`
	TemperatureUnit   string               `bigquery:"unit"`
	TemperatureValue  bigquery.NullFloat64 `bigquery:"temperature"`
	HeatSetPointValue bigquery.NullFloat64 `bigquery:"heat_setpoint"`
	HeatDemandValue   bigquery.NullFloat64 `bigquery:"heat_demand"`
	HumidityValue     bigquery.NullFloat64 `bigquery:"humidity"`
}
