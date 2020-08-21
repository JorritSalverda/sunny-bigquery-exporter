package main

import (
	"runtime"
	"time"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

var (
	// set when building the application
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()

	// application specific config
	sunnyHostIPAddress = kingpin.Flag("sunny-host-ip", "Host ip address of sunny inverter").Default("127.0.0.1").OverrideDefaultFromEnvar("SUNNY_HOST_IP").String()
	sunnyHostPort      = kingpin.Flag("sunny-host-port", "Host port of sunny inverter").Default("502").OverrideDefaultFromEnvar("SUNNY_HOST_PORT").Int()
	sunnyUnitID        = kingpin.Flag("sunny-unit-id", "ModBus unit id of sunny inverter").Default("3").OverrideDefaultFromEnvar("SUNNY_UNIT_ID").Int()

	bigqueryEnable    = kingpin.Flag("bigquery-enable", "Toggle to enable or disable bigquery integration").Default("true").OverrideDefaultFromEnvar("BQ_ENABLE").Bool()
	bigqueryProjectID = kingpin.Flag("bigquery-project-id", "Google Cloud project id that contains the BigQuery dataset").Envar("BQ_PROJECT_ID").Required().String()
	bigqueryDataset   = kingpin.Flag("bigquery-dataset", "Name of the BigQuery dataset").Envar("BQ_DATASET").Required().String()
	bigqueryTable     = kingpin.Flag("bigquery-table", "Name of the BigQuery table").Envar("BQ_TABLE").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate))

	// init bigquery client
	bigqueryClient, err := NewBigQueryClient(*bigqueryProjectID, *bigqueryEnable)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating bigquery client")
	}

	// init bigquery table if it doesn't exist yet
	initBigqueryTable(bigqueryClient)

	client, err := NewSunnyBoyClient(*sunnyHostIPAddress, *sunnyHostPort, *sunnyUnitID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating sunny boy client")
	}

	totalWhOut, err := client.GetTotalWhOut()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading totalWhOut")
	}

	measurement := BigQueryMeasurement{
		Readings: []BigQueryInverterReading{
			{
				Name:    "Sunny TriPower 8.0",
				Reading: float64(totalWhOut),
				Unit:    "Wh",
			},
		},
		InsertedAt: time.Now().UTC(),
	}

	err = bigqueryClient.InsertMeasurement(*bigqueryDataset, *bigqueryTable, measurement)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed inserting measurements into bigquery table")
	}
	log.Info().Msgf("Stored %v readings, exiting...", len(measurement.Readings))

	// done
	log.Info().Msg("Finished exporting metrics")
}

func initBigqueryTable(bigqueryClient BigQueryClient) {

	log.Debug().Msgf("Checking if table %v.%v.%v exists...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
	tableExist := bigqueryClient.CheckIfTableExists(*bigqueryDataset, *bigqueryTable)
	if !tableExist {
		log.Debug().Msgf("Creating table %v.%v.%v...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
		err := bigqueryClient.CreateTable(*bigqueryDataset, *bigqueryTable, BigQueryMeasurement{}, "inserted_at", true)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed creating bigquery table")
		}
	} else {
		log.Debug().Msgf("Trying to update table %v.%v.%v schema...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
		err := bigqueryClient.UpdateTableSchema(*bigqueryDataset, *bigqueryTable, BigQueryMeasurement{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed updating bigquery table schema")
		}
	}
}
