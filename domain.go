package main

import (
	"time"
)

type BigQueryMeasurement struct {
	Readings   []BigQueryInverterReading `bigquery:"readings"`
	InsertedAt time.Time                 `bigquery:"inserted_at"`
}

type BigQueryInverterReading struct {
	Name    string  `bigquery:"name"`
	Reading float64 `bigquery:"reading"`
	Unit    string  `bigquery:"unit"`
}
