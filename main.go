package main

import (
	stdlog "log"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// set when building the application
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()

	// application specific config
	sunnyBoyHostIPAddress = kingpin.Flag("sunny-boy-host-ip", "Host ip address of sunny boy").Default("127.0.0.1").OverrideDefaultFromEnvar("SUNNY_BOY_HOST_IP").String()
	sunnyBoyHostPort      = kingpin.Flag("sunny-boy-host-port", "Host port of sunny boy").Default("9522").OverrideDefaultFromEnvar("SUNNY_BOY_HOST_PORT").Int()
	sunnyBoyUser          = kingpin.Flag("sunny-boy-user", "Username to log in to sunny boy").Default("User").OverrideDefaultFromEnvar("SUNNY_BOY_USER").String()
	sunnyBoyPassword      = kingpin.Flag("sunny-boy-password", "Password to log in to sunny boy").Default("0000").OverrideDefaultFromEnvar("SUNNY_BOY_PASSWORD").String()

	udpLocalPort = kingpin.Flag("udp-local-port", "Local port used for udp listener").Default("8855").OverrideDefaultFromEnvar("UDP_LOCAL_PORT").Int()

	bigqueryProjectID = kingpin.Flag("bigquery-project-id", "Google Cloud project id that contains the BigQuery dataset").Envar("BQ_PROJECT_ID").Required().String()
	bigqueryDataset   = kingpin.Flag("bigquery-dataset", "Name of the BigQuery dataset").Envar("BQ_DATASET").Required().String()
	bigqueryTable     = kingpin.Flag("bigquery-table", "Name of the BigQuery table").Envar("BQ_TABLE").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log as severity for stackdriver logging to recognize the level
	zerolog.LevelFieldName = "severity"

	// set some default fields added to all logs
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", app).
		Str("version", version).
		Logger()

	// use zerolog for any logs sent via standard log library
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// log startup message
	log.Info().
		Str("branch", branch).
		Str("revision", revision).
		Str("buildDate", buildDate).
		Str("goVersion", goVersion).
		Msgf("Starting %v version %v...", app, version)

	// bigqueryClient, err := NewBigQueryClient(*bigqueryProjectID)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed creating bigquery client")
	// }
	// initBigqueryTable(bigqueryClient)

	client, err := NewSunnyBoyClient(*udpLocalPort, *sunnyBoyHostIPAddress, *sunnyBoyHostPort, *sunnyBoyUser, *sunnyBoyPassword)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating sunny boy client")
	}

	results, err := client.ReadInputRegisters(8, 1)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading input registers")
	}

	log.Info().Interface("results", results).Msg("Retrieved results from reading input registers")

	// log.Debug().Msgf("Inserting measurements into table %v.%v.%v...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
	// err = bigqueryClient.InsertMeasurements(*bigqueryDataset, *bigqueryTable, measurements)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed inserting measurements into bigquery table")
	// }

	// done
	log.Info().Msg("Finished exporting metrics")
}

func initBigqueryTable(bigqueryClient BigQueryClient) {

	log.Debug().Msgf("Checking if table %v.%v.%v exists...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
	tableExist := bigqueryClient.CheckIfTableExists(*bigqueryDataset, *bigqueryTable)
	if !tableExist {
		log.Debug().Msgf("Creating table %v.%v.%v...", *bigqueryProjectID, *bigqueryDataset, *bigqueryTable)
		err := bigqueryClient.CreateTable(*bigqueryDataset, *bigqueryTable, BigQueryMeasurement{}, "measured_at", true)
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
