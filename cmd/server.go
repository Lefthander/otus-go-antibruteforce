package cmd
import (
	"flag"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/Lefthander/otus-go-antibruteforce/db"
	"github.com/Lefthander/otus-go-antibruteforce/logger"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "abf-srv",
	Short:"abf-srv start",
	Long:"abf-srv start",
	Run: func(cmd *cobra.Command, args []string) {
		configLocation:=flag.String("config","config.json","configuration file")
		flag.Parse()
		cfg, err:= config.GetConfig(*configLocation)
		
		if err != nil {
			log.Fatal("Error cannot get a config file ",err)
		}

		lg,err := logger.GetLogger(cfg)

		if err != nil {
			log.Fatal("Error cannot setup a logger",err)
		}

		psql,err := db.ConnectDB(cfg)
		
		if err != nil {
			log.Fatal("Error to setup a Postgresql",err)
		}
		

	}


}