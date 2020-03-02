package main

import (
	"flag"
	"log"
	"net"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/Lefthander/otus-go-antibruteforce/db"
	"github.com/Lefthander/otus-go-antibruteforce/internal/adapters"
	ipdbs "github.com/Lefthander/otus-go-antibruteforce/internal/adapters/db"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/interfaces"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/usecases"
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc"
	"github.com/Lefthander/otus-go-antibruteforce/internal/metrics"
	"github.com/Lefthander/otus-go-antibruteforce/logger"
	"github.com/spf13/cobra"
)

// RootCmd is a main command to run the service
var RootCmd = &cobra.Command{ //nolint
	Use:   "abf-srv",
	Short: "abf-srv to Run the ABF grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		var store interfaces.FilterKeeper
		configLocation := flag.String("config", "config.yaml", "configuration file")
		flag.Parse()

		err := config.GetConfig(*configLocation)

		if err != nil {
			log.Fatal("Error cannot get a config file ", err)
		}

		loggercfg := config.GetLoggerCfg()

		if loggercfg == nil {
			log.Fatal("Failed to get logger config")
		}

		lg, err := logger.GetLogger(loggercfg)

		if err != nil {
			log.Fatal("Error cannot setup a Zap logger", err)
		}

		dbcfg := config.GetDBCfg()

		if loggercfg.Verbose {
			lg.Info("Setting UP the B/W IP List store...")
		}

		switch dbcfg.DBType {
		case "memory":
			lg.Info("Selected Memory Map as storage for IP Filters")
			store = adapters.NewIPFilterMemory()
		case "psql":
			lg.Info("Selected Postgresql DB as storage for IP Filters")

			psql, err := db.ConnectDB(dbcfg)

			if err != nil {
				lg.Fatal("Error to setup a Postgresql connection", err)
			}
			store = ipdbs.NewDBStore(psql)

		default:
			lg.Info("Selected Memory Map as storage for IP Filters")
			store = adapters.NewIPFilterMemory()
		}

		if loggercfg.Verbose {
			lg.Info("Setting UP the Login/Password/IP Token Buckets...")
		}

		loginbucket := adapters.NewTokenBucketMemory()
		passwdbucket := adapters.NewTokenBucketMemory()
		ipbucket := adapters.NewTokenBucketMemory()

		cfg := config.GetServiceCfg()

		if loggercfg.Verbose {
			lg.Info("Setting UP the ABF Service...")
		}

		abfservice := usecases.NewABFService(cfg.ConstraintN, cfg.ConstraintM,
			cfg.ConstraintK, loginbucket, passwdbucket, ipbucket, store, lg, cfg)

		abfserver := grpc.NewABFServer(abfservice)

		metr := &metrics.MonService{
			Port: cfg.MonitorPort,
		}

		if loggercfg.Verbose {
			lg.Info("Starting Prometheus metric handler...")
		}

		metr.ServeMetrics()

		if loggercfg.Verbose {
			lg.Infof("Starting AntiBruteForceService on Port=%s...", cfg.ServicePort)
		}
		err = abfserver.ListenAndServe(net.JoinHostPort("", cfg.ServicePort))

		if err != nil {
			lg.Fatal("Error cannot start AntiBruteForce Server", err)
		}
	},
}

func init() { // nolint
}

func main() {
	log.Println("Starting ABF...")

	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Failed to start...", err)
	}
}
