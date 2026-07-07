package main

import (
	"fmt"

	"github.com/MonarchRyuzaki/dns-resolver-go/resolver"
)

func main() {
	queryBytes := resolver.NewSimpleQuery(22, "dns.google.com")
	fmt.Printf("Bytes to send: %x\n", queryBytes)
}
