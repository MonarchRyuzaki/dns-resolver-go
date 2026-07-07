package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

func parseFlags() string {
	domain := flag.String("domain", "dns.google.com", "The domain to resolve")
	flag.Parse()
	return *domain
}

// resolve queries a specific DNS server for the given domain and returns the parsed DNSMessage
func resolve(domain string, serverIP string) resolver.DNSMessage {
	queryBytes := resolver.NewSimpleQuery(22, domain)

	conn, err := net.Dial("udp", serverIP+":53")
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", serverIP, err)
	}
	defer conn.Close()

	_, err = conn.Write(queryBytes)
	if err != nil {
		log.Fatalf("Failed to send query: %v", err)
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	return resolver.DecodeResponse(buf[:n])
}

func main() {
	domain := parseFlags()
	fmt.Printf("Resolving %s...\n", domain)

	msg := resolve(domain, "8.8.8.8")

	for _, ans := range msg.Answers {
		// Only print A Records (Type 1)
		if ans.Type == 1 {
			fmt.Printf("Resolved IP: %v\n", ans.IPString())
		}
	}
}
