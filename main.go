package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

type Options struct {
	Domain    string
	Recursive bool
}

func parseFlags() Options {
	domain := flag.String("domain", "dns.google.com", "The domain to resolve")
	recursive := flag.Bool("recursive", false, "Use recursive resolution (RD=1) via Google DNS")
	flag.Parse()
	return Options{Domain: *domain, Recursive: *recursive}
}

// resolve queries a specific DNS server for the given domain and returns the parsed DNSMessage
func resolve(domain string, serverIP string, recursive bool) resolver.DNSMessage {
	queryBytes := resolver.NewSimpleQuery(22, domain, recursive)

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
	opts := parseFlags()

	serverIP := "198.41.0.4"
	if opts.Recursive {
		fmt.Printf("Resolving %s recursively...\n", opts.Domain)
		serverIP = "8.8.8.8"
	} else {
		fmt.Printf("Resolving %s iteratively...\n", opts.Domain)
	}

	for {
		fmt.Printf("Querying %s for %s\n", serverIP, opts.Domain)
		msg := resolve(opts.Domain, serverIP, opts.Recursive)

		foundAnswer := false
		for _, ans := range msg.Answers {
			if ans.Type == 1 {
				fmt.Printf("Resolved IP: %v\n", ans.IPString())
				foundAnswer = true
			}
		}

		if foundAnswer {
			return
		}

		foundGlue := false
		for _, add := range msg.Additionals {
			if add.Type == 1 {
				serverIP = add.IPString()
				foundGlue = true
				break
			}
		}

		if foundGlue {
			continue
		}

		if len(msg.Authorities) > 0 {
			// This means they gave us the name of the next server, but didn't give us the IP in Additionals.
			// A true production resolver would now recursively call `resolve(authorityName, "198.41.0.4")`
			// to find the IP of the nameserver, and then continue.
			log.Fatalf("Received NS referral without a glue IP! Full nameserver resolution not yet implemented.")
		}

		log.Fatalf("Resolution failed: No answers, no glue records, and no authorities returned.")
	}
}
