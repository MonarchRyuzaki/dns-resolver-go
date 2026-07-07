package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

func main() {
	queryBytes := resolver.NewSimpleQuery(22, "dns.google.com")
	// fmt.Printf("Bytes to send: %x\n", queryBytes)

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

	// 5. Pass the valid bytes to your decoder
	msg := resolver.DecodeResponse(buf[:n])
	fmt.Printf("Resolved IP: %v\n", msg.Answer.IPString())
}
