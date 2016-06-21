package main

import (
	"reflect"
	"runtime"
	"time"
)

func PeriodicJob(interval int, f func(*app_state, *app_config, *traffic_data), AppState *app_state, AppConfig *app_config, TrafficData *traffic_data) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	funcname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	AppState.Wait.Add(1)
	go func() {
		DebugLogger.Printf("Starting periodic job %s every %d seconds\n", funcname, interval)
		defer AppState.Wait.Done()
		for AppState.Running {
			select {
			case <-ticker.C:
				f(AppState, AppConfig, TrafficData)
			}
		}
		DebugLogger.Printf("Finished periodic job %s\n", funcname)
	}()
}
