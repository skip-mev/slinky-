package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	oracle2 "github.com/CosmWasm/wasmd/service/servers/oracle"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof" //nolint: gosec

	"go.uber.org/zap"
)

var (
	host    = flag.String("host", "0.0.0.0", "host for the grpc-service to listen on")
	port    = flag.String("port", "8080", "port for the grpc-service to listen on")
	dataHex = flag.String("datahex", "0xdeadbeef", "vote extension data to return")
)

// start the oracle-grpc server + oracle process, cancel on interrupt or terminate.
func main() {
	// channel with width for either signal
	sigs := make(chan os.Signal, 1)

	// gracefully trigger close on interrupt or terminate signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// parse flags
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %s\n", err.Error())
		return
	}

	databytes, err := hex.DecodeString(*dataHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode vote extension data: %s\n", err.Error())
		return
	}

	fmt.Println("databytes: ", string(databytes))

	// create server
	conn, err := ethclient.Dial("https://ethereum-rpc.publicnode.com")
	if err != nil {
		return
	}
	srv := oracle2.NewOracleServer(logger, databytes, conn)

	// cancel oracle on interrupt or terminate
	go func() {
		<-sigs
		logger.Info("received interrupt or terminate signal; closing oracle")

		cancel()
	}()

	// start oracle + server, and wait for either to finish
	if err := srv.StartServer(ctx, *host, *port); err != nil {
		logger.Error("stopping server", zap.Error(err))
	}
}
