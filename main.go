package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

func main() {
	domain := flag.String("domain", "dns.google.com", "The domain to resolve")
	flag.Parse()

	queryBytes := resolver.NewSimpleQuery(22, *domain)
	fmt.Printf("Resolving %s...\n", *domain)

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

	msg := resolver.DecodeResponse(buf[:n])

	for _, ans := range msg.Answers {
		if ans.Type == 1 {
			fmt.Printf("Resolved IP: %v\n", ans.IPString())
		}
	}
}
