package main

import (
	"time"
)

type traffic_point struct {
	BytesPerSecond      int
	AllPacketsPerSecond int
	UDPPacketsPerSecond int
	SYNPacketsPerSecond int
}

type host_traffic struct {
	TrrafficPoints [60]traffic_point
}

type traffic_data struct {
	HostTrraffic map[[16]byte]host_traffic
}

func _CountersRotator(AppState *app_state, TrafficData *traffic_data) {
	DebugLogger.Println("Rotating stuff")
}

func CountersRotator(AppState *app_state, TrafficData *traffic_data) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		DebugLogger.Println("CounterRotator starting")
		defer AppState.Wait.Done()
		for AppState.Running {
			select {
			case <-ticker.C:
				_CountersRotator(AppState, TrafficData)
			}
		}
		DebugLogger.Println("CounterRotator finished")
	}()
}
