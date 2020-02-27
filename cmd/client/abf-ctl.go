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
	client api.ABFServiceClient //nolint

	clientcfg *config.ClientConfig // nolint
)

// RootCmd is a main command to handle the client commands
var RootCmd = &cobra.Command{ // nolint
	Use:       "abf-ctl [add-to-iplist,delete-from-iplist,reset-limits,show-iplist,test",
	Short:     "abf-ctl gRPC client for AntiBruteForce Service",
	ValidArgs: []string{"add-to-iplist", "delete-from-iplist", "reset-limits", "show-iplist", "test"},
	Args:      cobra.ExactValidArgs(numberOfValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		clientcfg = config.GetClientCfg()
		ctx, cancel := context.WithTimeout(context.Background(), clientcfg.ConnectionTimeOut)
		client = newClient(ctx, clientcfg.Host, clientcfg.Port)

		// Running watchdog goroutine to control the system interrupts
		go func() {
			terminate := make(chan os.Signal, 1)
			signal.Notify(terminate, os.Interrupt, syscall.SIGINT)
			<-terminate
			log.Println("Received system interrupt...")
			cancel()
		}()

		switch args[0] {
		case "add-to-iplist":
			addToIPList(ctx)
		case "delete-from-iplist":
			delFromIPList(ctx)
		case "reset-limits":
			resetLimits(ctx)
		case "show-iplist":
			showIPList(ctx)
		case "test":
			isConform(ctx)
		}

	},
}

func newClient(ctx context.Context, host, port string) api.ABFServiceClient {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal("Cannot connect to ABF server", err)
	}

	return api.NewABFServiceClient(conn)

}

func main() {
	log.Println("ABF Client...")

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func isConform(ctx context.Context) {

}

func addToIPList(ctx context.Context) {

}

func delFromIPList(ctx context.Context) {

}

func resetLimits(ctx context.Context) {

}

func showIPList(ctx context.Context) {

}
