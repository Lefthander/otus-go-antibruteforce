package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lefthander/otus-go-antibruteforce/config"
	"github.com/Lefthander/otus-go-antibruteforce/internal/grpc/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	numberOfValidArgs = 1
)

var (
	//	client api.ABFServiceClient //nolint

	//	clientcfg *config.ClientConfig // nolint

	logcfg *config.LoggerConfig // nolint

	login    string //nolint
	network  string // nolint
	password string // nolint

	color bool // nolint

	ipaddress string //nolint

	//ctx context.Context // nolint
)

// RootCmd is a main command to handle the client commands
var RootCmd = &cobra.Command{ // nolint
	Use:       "abf-ctl",
	Short:     "abf-ctl gRPC client for ABF Service",
	ValidArgs: []string{"add", "del", "reset", "show", "test"},
	Args:      cobra.ExactValidArgs(numberOfValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var addCmd = &cobra.Command{ //nolint
	Use:   "add",
	Short: "add to black/white list",
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg := config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		defer cancel()

		client := newClient(ctx, clientcfg.Host, clientcfg.Port)

		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()

		r, err := client.AddToIpFilter(ctx, &api.IPFilterData{Network: network, Color: color})

		if err != nil {
			log.Fatalf("unable to add to list: %v", err)
		}

		log.Println("Done: ", r.Error)
	},
}

var delCmd = &cobra.Command{ //nolint
	Use:   "del",
	Short: "del from black/white list",
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg := config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		defer cancel()

		client := newClient(ctx, clientcfg.Host, clientcfg.Port)

		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()
		r, err := client.DeleteFromIpFilter(ctx, &api.IPFilterData{Network: network, Color: color})

		if err != nil {
			log.Fatalf("unable to add to list: %v", err)
		}

		log.Println("Done: ", r.Error)
	},
}
var resetCmd = &cobra.Command{ //nolint
	Use:   "reset",
	Short: "reset",
	Long:  "reset the bucket limits for specifi pair of Login/IP",
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg := config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		defer cancel()

		client := newClient(ctx, clientcfg.Host, clientcfg.Port)

		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()

		r, err := client.Reset(ctx, &api.AuthRequest{Login: login, Password: "", Ipaddr: ipaddress})

		if err != nil {
			log.Fatalf("unable to reset limits: %v", err)
		}

		log.Println("Done: ", r.Response)
	},
}
var showCmd = &cobra.Command{ //nolint
	Use:   "show",
	Short: "show",
	Long:  "show dumps the corresponding ip table black/white",
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg := config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		defer cancel()

		client := newClient(ctx, clientcfg.Host, clientcfg.Port)

		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()

		r, err := client.GetIpFilters(ctx, &api.IPFilterData{Network: "", Color: color})

		if err != nil {
			log.Fatalf("unable to show the ip filters: %v", err)
		}
		log.Println("Done: ", r.Filters)
	},
}

var testCmd = &cobra.Command{ //nolint
	Use:   "test",
	Short: "test",
	Long:  "test request ABF service to verify the triplet of (login,password,ip)",
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg := config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		defer cancel()

		client := newClient(ctx, clientcfg.Host, clientcfg.Port)

		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()

		r, err := client.Allow(ctx, &api.AuthRequest{Login: login, Password: password, Ipaddr: ipaddress})

		if err != nil {
			log.Fatalf("unable to show the ip filters: %v", err)
		}
		log.Println("Done: ", r.GetOk())
	},
}

func newClient(ctx context.Context, host, port string) api.ABFServiceClient {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(host, port), grpc.WithInsecure())

	if err != nil {
		log.Fatal("Cannot connect to ABF server", err)
	}

	return api.NewABFServiceClient(conn)
}

func init() { // nolint
	addCmd.PersistentFlags().StringVarP(&network, "network", "n", "", "network")
	addCmd.PersistentFlags().BoolVarP(&color, "color", "c", true, "white - true, false black")
	RootCmd.AddCommand(addCmd)
	delCmd.PersistentFlags().StringVarP(&network, "network", "n", "", "network")
	delCmd.PersistentFlags().BoolVarP(&color, "color", "c", true, "white - true, false black")
	RootCmd.AddCommand(delCmd)
	resetCmd.PersistentFlags().StringVarP(&login, "login", "l", "", "login to reset the limits")
	resetCmd.PersistentFlags().StringVarP(&ipaddress, "ipaddress", "i", "", "ip to reset to reset the limits")
	RootCmd.AddCommand(resetCmd)
	showCmd.PersistentFlags().BoolVarP(&color, "color", "c", true, "white - true, false black")
	RootCmd.AddCommand(showCmd)
	testCmd.PersistentFlags().StringVarP(&login, "login", "l", "", "login to test")
	testCmd.PersistentFlags().StringVarP(&login, "password", "p", "", "password to test")
	testCmd.PersistentFlags().StringVarP(&ipaddress, "ipaddress", "i", "", "ip to test")
	RootCmd.AddCommand(testCmd)
}

func main() {
	log.Println("ABF Client started..")

	err := config.GetConfig("config.yml")

	if err != nil {
		log.Fatal("Error setting up the configuration ", err)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
