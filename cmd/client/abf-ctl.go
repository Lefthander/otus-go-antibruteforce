package main

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	ConnectTimeOut = 30 * time.Second
)
var RootCmd = &cobra.Command{
Use:"abf-ctl [add-to-iplist,delete-from-iplist,reset,show-iplist",
Short:"abf-ctl gRPC client for AntiBruteForce Service",
ValidArgs:[]string{"add-to-iplist","delete-from-iplist","reset","show-iplist"},
Args: cobra.ExactValidArgs(1),
Run: func(cmd *cobra.Command,args []string) {
	ctx,cancel:= context.WithTimeOut(context.Background(),ConnectTimeOut)
	}
},
}