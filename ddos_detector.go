package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type app_state struct {
	Wait    sync.WaitGroup
	Running bool
}

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func InitLogging(WithDebug bool) {
	InfoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	var DebugWriter io.Writer
	if WithDebug {
		DebugWriter = os.Stdout
	} else {
		DebugWriter = ioutil.Discard
	}
	DebugLogger = log.New(DebugWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

}

func KillSignal(Running *bool) {
	SignalChannel := make(chan os.Signal, 1)
	signal.Notify(SignalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-SignalChannel
		InfoLogger.Println("Got SIGTERM, finishing work gracefully.")
		*Running = false
		signal.Reset() // Next ctrl+c will effect in ungraceful stop
	}()

}

func MainLoop(AppState *app_state, AppConfig *app_config, TrafficData *traffic_data) {
	AppState.Running = true

	AppState.Wait.Add(1)
	go sFlowListener(AppState, AppConfig.SFlowConfig)

	AppState.Wait.Add(1)
	go CountersRotator(AppState, TrafficData)

	AppState.Wait.Wait()
}

func main() {
	var (
		configFile  string
		withDebug   bool
		AppState    app_state
		AppConfig   app_config
		TrafficData traffic_data
		err         error
	)

	flag.StringVar(&configFile, "c", "/etc/ddos_detector.toml", "Path to configuration file.")
	flag.BoolVar(&withDebug, "v", false, "Be verbose, show debugging output.")
	flag.Parse()

	InitLogging(withDebug)
	InfoLogger.Println("Starting DDoS Detector.")

	AppConfig, err = ReadConfig(configFile)
	if err != nil {
		os.Exit(1)
	}

	KillSignal(&AppState.Running)
	MainLoop(&AppState, &AppConfig, &TrafficData)

	InfoLogger.Println("DDoS Detector finished.")
}
