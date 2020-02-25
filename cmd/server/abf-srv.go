package main

import (
	"flag"
	"log"
	"net"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/Lefthander/otus-go-antibruteforce/db"
	"github.com/Lefthander/otus-go-antibruteforce/internal/adapters"
	ipdbs "github.com/Lefthander/otus-go-antibruteforce/internal/adapters/db"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/usecases"
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc"
	"github.com/Lefthander/otus-go-antibruteforce/internal/metrics"
	"github.com/Lefthander/otus-go-antibruteforce/logger"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "abf-srv",
	Short: "abf-srv to Run the ABF grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		configLocation := flag.String("config", "config.json", "configuration file")
		flag.Parse()

		cfg, err := config.GetConfig(*configLocation)

		if err != nil {
			log.Fatal("Error cannot get a config file ", err)
		}

		lg, err := logger.GetLogger(cfg)

		if err != nil {
			log.Fatal("Error cannot setup a Zap logger", err)
		}

		psql, err := db.ConnectDB(cfg)

		if err != nil {
			log.Fatal("Error to setup a Postgresql connection", err)
		}

		ipstore := ipdbs.NewDBStore(psql)

		loginbucket := adapters.NewTokenBucketMemory()
		passwdbucket := adapters.NewTokenBucketMemory()
		ipbucket := adapters.NewTokenBucketMemory()

		abfservice := usecases.NewABFService(cfg.ConstraintN, cfg.ConstraintM, cfg.ConstraintK, loginbucket, passwdbucket, ipbucket, ipstore, lg, cfg)

		abfserver := grpc.NewABFServer(abfservice)

		metr := &metrics.MonService{
			Port: cfg.MonitorPort,
		}

		log.Println("Starting Prometheus metric handler...")

		metr.ServeMetrics()

		log.Println("Starting AntiBruteForceService...")

		err = abfserver.ListenAndServe(net.JoinHostPort(cfg.ServiceHost, cfg.ServicePort))

		if err != nil {
			log.Fatal("Error cannot start AntiBruteForce Server", err)
		}
	},
}

func init() {

}

func main() {
	log.Println("Starting ABF...")

	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Failed to start...", err)
	}
}
