package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	grpc2 "github.com/nitrictech/nitric/core/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var runtime string
var app string
var timeout int

var rootCmd = &cobra.Command{
	Use:   "nitricd",
	Short: "nitricd - Nitric daemon for custom runtimes",
	Long: `nitricd
   
the nitric daemon process facilitates communication between nitric workers and a custom gateway`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running root command")
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		cobra.CheckErr(err)

		faasAddress := lis.Addr().String()
		faasAddressEnv := fmt.Sprintf("NITRIC_FAAS_ADDRESS=%s", faasAddress)

		fmt.Printf("faas gateway listening on %s\n", faasAddress)

		srv := grpc.NewServer()

		minWorkers := utils.GetEnv("MIN_WORKERS", "0")
		mw, err := strconv.Atoi(minWorkers)
		cobra.CheckErr(err)

		pool := pool.NewProcessPool(&pool.ProcessPoolOptions{
			MinWorkers: mw,
		})

		faasSrv := grpc2.NewFaasServer(pool)
		// Broker needs a FaaS server
		v1.RegisterFaasServiceServer(srv, faasSrv)
		// Broker also provider unary layer for trigger workers (this is for custom gateways)
		v1.RegisterFaasGatewayServiceServer(srv, faasSrv)

		// Start the gRPC service
		// TODO: Handle errors...
		fmt.Println("Start FaaS services")
		go srv.Serve(lis)

		// cobra.CheckErr(err)
		// Start the nitric application
		// execute the provided script for the --app flag
		fmt.Println("Starting app")
		appArgs := strings.Split(app, " ")
		appCmd := exec.Command(appArgs[0], appArgs[1:]...)
		appCmd.Env = append(os.Environ(), faasAddressEnv)
		appCmd.Stdout = os.Stdout

		// Start the application non-blocking
		err = appCmd.Start()
		cobra.CheckErr(err)

		// Once the workers are available (up to min workers) we can kick off the runtime (which will expose the gateway)
		fmt.Println("waiting for workers")
		err = pool.WaitForMinimumWorkers(timeout)
		cobra.CheckErr(err)

		fmt.Println("starting runtime")
		runtimeArgs := strings.Split(runtime, " ")
		runtimeCmd := exec.Command(runtimeArgs[0], runtimeArgs[1:]...)
		runtimeCmd.Env = append(runtimeCmd.Env, faasAddressEnv)
		runtimeCmd.Stdout = os.Stdout

		err = runtimeCmd.Run()
		cobra.CheckErr(err)
	},
}

// Start a new broker
func main() {
	// the runtime flag provides direction on starting the nitric runtime process
	rootCmd.Flags().StringVarP(&runtime, "runtime", "r", "", "command to start runtime process")
	rootCmd.MarkFlagRequired("runtime")
	// the app flag provides direction on starting the nitric application process
	rootCmd.Flags().StringVarP(&app, "app", "a", "", "command to start application process")
	rootCmd.MarkFlagRequired("app")

	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "timeout to wait for workers to start")

	err := rootCmd.Execute()

	cobra.CheckErr(err)
}
