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

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	wait        sync.WaitGroup
	running     bool
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

func KillSignal() {
	SignalChannel := make(chan os.Signal, 1)
	signal.Notify(SignalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-SignalChannel
		InfoLogger.Println("Got SIGTERM, finishing work gracefully.")
		running = false
		signal.Reset() // Next ctrl+c will effect in ungraceful stop
	}()

}

func MainLoop(AppConfig app_config) {
	running = true

	wait.Add(1)
	go sFlowListener(AppConfig)

	wait.Wait()
}

func main() {
	var configFile string
	var withDebug bool
	var AppConfig app_config

	flag.StringVar(&configFile, "c", "/etc/ddos_detector.toml", "Path to configuration file.")
	flag.BoolVar(&withDebug, "v", false, "Be verbose, show debugging output.")
	flag.Parse()

	InitLogging(withDebug)
	InfoLogger.Println("Starting DDoS Detector.")

	if !ReadConfig(configFile, &AppConfig) {
		os.Exit(1)
	}

	KillSignal()
	MainLoop(AppConfig)

	InfoLogger.Println("DDoS Detector finished.")
}
