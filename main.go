package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sqooba/go-common/logging"
	"github.com/sqooba/go-common/version"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	setLogLevel = flag.String("set-log-level", "", "Change log level. Possible values are trace,debug,info,warn,error,fatal,panic")
	log         = logging.NewLogger()
)

type envConfig struct {
	Port             string `envconfig:"PORT" default:"8080"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
	MetricsNamespace string `envconfig:"METRICS_NAMESPACE" default:"autoheal"`
	MetricsSubsystem string `envconfig:"METRICS_SUBSYSTEM" default:""`
	MetricsPath      string `envconfig:"METRICS_PATH" default:"/metrics"`

	ContainerLabel     string        `envconfig:"AUTOHEAL_CONTAINER_LABEL" default:"autoheal"`
	Interval           time.Duration `envconfig:"AUTOHEAL_INTERVAL" default:"5s"`
	StartPeriod        time.Duration `envconfig:"AUTOHEAL_START_PERIOD" default:"0s"`
	DefaultStopTimeout time.Duration `envconfig:"AUTOHEAL_DEFAULT_STOP_TIMEOUT" default:"10s"`
	CurlTimeout        time.Duration `envconfig:"CURL_TIMEOUT" default:"30s"`
}

func main() {
	log.Println("docker-autoheal application is starting...")
	log.Printf("Version    : %s", version.Version)
	log.Printf("Commit     : %s", version.GitCommit)
	log.Printf("Build date : %s", version.BuildDate)
	log.Printf("OSarch     : %s", version.OsArch)

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s\n", err)
		return
	}

	flag.Parse()
	err := logging.SetLogLevel(log, env.LogLevel)
	if err != nil {
		log.Fatalf("Logging level %s do not seem to be right. Err = %v", env.LogLevel, err)
	}

	if *setLogLevel != "" {
		logging.SetRemoteLogLevelAndExit(log, env.Port, *setLogLevel)
	}

	// Special endpoint to change the verbosity at runtime, i.e. curl -X PUT --data debug ...
	logging.InitVerbosityHandler(log, http.DefaultServeMux)
	initPrometheus(env, http.DefaultServeMux)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	dockerHandler, err := newDockerHandler(env)
	if err != nil {
		log.Fatalf("Got an exception while initialising docker handler. Err = %v", err)
	}
	go dockerHandler.handle()

	s := http.Server{Addr: fmt.Sprint(":", env.Port)}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	<-signalChan
	log.Printf("Shutdown signal received, exiting...")

	err = s.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Got an error while shutting down: %v\n", err)
	}

	// Wait for processing to complete properly
	wg.Wait()
}
