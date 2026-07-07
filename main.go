package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

func main() {
	queryBytes := resolver.NewSimpleQuery(22, "dns.google.com")
	fmt.Printf("Bytes to send: %x\n", queryBytes)

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err.Error())
	}

	conn.Write(queryBytes)

	buf := make([]byte, 512)

	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err.Error())
	}

	record := resolver.DecodeResponse(buf[:n])
	fmt.Printf("Resolved IP: %v\n", record.IPString())
}
