package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/nats-rider/api-server"
	"github.com/wallyqs/nats-rider/kit/component"
)

func main() {
	var (
		showHelp     bool
		showVersion  bool
		serverListen string
		natsServers  string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: api-server [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&serverListen, "l", "0.0.0.0:9090", "Network host:port to listen on")
	flag.StringVar(&serverListen, "listen", "0.0.0.0:9090", "Network host:port to listen on")
	flag.StringVar(&natsServers, "s", nats.DefaultURL, "List of NATS Servers to connect")
	flag.StringVar(&natsServers, "servers", nats.DefaultURL, "List of NATS Servers to connect")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATS Rider API Server v%s\n", apiserver.Version)
		os.Exit(0)
	}
	log.Printf("Starting NATS Rider API Server version %s", apiserver.Version)

	// Register new component within the system.
	comp := kit.NewComponent("api-server")

	// Connect to NATS and setup discovery subscriptions.
	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}
	s := apiserver.Server{
		Component: comp,
	}

	err = s.ListenAndServe(serverListen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for HTTP requests on %v", serverListen)

	// TODO: Change into using context for shutdown and cancellation
	// TODO: Handle signals!
	runtime.Goexit()
}
