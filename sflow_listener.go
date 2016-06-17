package main

import (
	"bytes"
	"github.com/Cistern/sflow"
	"net"
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

func sFlowListener() (err error) {
	// Start listening UDP socket, check if it started properly
	UDPAddr, err := net.ResolveUDPAddr("udp", ":6343")
	conn, err := net.ListenUDP("udp", UDPAddr)
	if err != nil {
		return err
	}

	var buffer []byte
	for {
		conn.ReadFromUDP(buffer)
		sFlowParser(buffer)
	}

	return nil
}
