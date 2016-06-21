package main

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

func CountersRotator(AppState *app_state, AppConfig *app_config, TrafficData *traffic_data) {
	DebugLogger.Println("Rotating stuff")
}
