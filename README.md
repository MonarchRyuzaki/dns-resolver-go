# DNS Resolver (Go)

A fully functioning DNS resolver built entirely from scratch in Go. This project implements the core specifications defined in [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035), manually crafting UDP packets, manipulating bitwise headers, and decoding DNS compression pointers without relying on any external DNS libraries.

## Features

- **Manual Packet Serialization:** Implements custom binary packing and unpacking for DNS Headers, Questions, and Records.
- **DNS Compression Handling:** Parses and unwraps 16-bit compression pointers (the notorious `0xC0` trap) to successfully reconstruct dynamically compressed domain names.
- **Iterative Resolution (Default):** Acts as a true infrastructure-level DNS resolver by starting at the Root Servers (`198.41.0.4`) and manually traversing the DNS tree (Root -> TLD -> Authoritative) using Glue Records until the final IP is found.
- **Recursive Resolution:** Supports offloading the resolution work to a public resolver (like Google's `8.8.8.8`) via a CLI flag.

## Usage

Run the resolver directly from the terminal using the Go CLI.

**Iterative Resolution (Default):**
Walks the DNS tree starting from the Root Servers.
```bash
go run main.go -domain=github.com
```

**Recursive Resolution:**
Sets the `RD=1` (Recursion Desired) flag and queries Google's DNS to do the heavy lifting.
```bash
go run main.go -domain=github.com -recursive
```

## Project Structure
- `main.go`: The CLI entry point and the core Iterative Resolution loop.
- `resolver/types.go`: Struct definitions for the DNS Message anatomy (Header, Question, Record).
- `resolver/encoder.go`: Bitwise logic for translating Go structs into raw UDP byte slices.
- `resolver/decoder.go`: Logic for parsing raw UDP responses back into structs, including the complex `DecodeDomainName` algorithm.
- `docs/`: Contains technical write-ups and explanations on complex DNS concepts (like Compression).

## Future Roadmap

While the core IPv4 (A Record) resolution is fully operational, there are several features planned for future implementation:

- **CNAME Support:** Add parsing and follow-through logic for Canonical Name (CNAME) aliases (Type 5).
- **Robust Authority Resolution:** Currently, the iterative loop relies on IP addresses being provided in the `Additionals` (Glue) section. If an Authoritative server is returned by name only, the resolver needs to pause, recursively resolve that Nameserver's IP from scratch, and then continue.
- **Local Caching:** Implement a local, TTL-respecting memory cache to store resolved IPs and TLD Nameserver IPs. This will significantly speed up subsequent queries and reduce network load on the Root Servers.
