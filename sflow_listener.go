package main

import (
	"bytes"
	"fmt"
	"github.com/Cistern/sflow"
	"net"
	"time"
)

func sFlowParser(buffer []byte) {
	reader := bytes.NewReader(buffer)
	d := sflow.NewDecoder(reader)
	dgram, err := d.Decode()
	if err != nil {
		println(err)
		return
	}
	for _, sample := range dgram.Samples {
		println(sample)
	}
}

func sFlowListener(AppConfig app_config) (err error) {
	defer wait.Done()

	var udp_addr = fmt.Sprintf("[%s]:%d", AppConfig.SFlowConfig.Address, AppConfig.SFlowConfig.Port)

	DebugLogger.Println("Binding sFlow listener to", udp_addr)

	UDPAddr, err := net.ResolveUDPAddr("udp", udp_addr)
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	conn, err := net.ListenUDP("udp", UDPAddr)
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}

	var buffer []byte
	for running {
		/*
		  Normally read would block, but we want to be able to break this
		  loop gracefuly. So add read timeout and every 0.1s check if it is
		  time to finish
		*/
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		var read, _, err = conn.ReadFromUDP(buffer)
		if read > 0 && err != nil {
			sFlowParser(buffer)
		}

	}

	conn.Close()

	return nil
}
