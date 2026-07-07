# Custom DNS Resolver (Go)

A fully functioning DNS resolver built entirely from scratch in Go, interacting directly with raw UDP sockets and implementing the core specifications of [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035).

Rather than relying on external libraries, this project prioritizes a solid grasp of low-level networking primitives, featuring manual bitwise UDP packet crafting, dynamic DNS pointer decompression, and full iterative DNS tree traversal.

## Core Architecture

*   **Iterative Resolution Engine:** Acts as a true infrastructure-level DNS resolver. Drops the Recursion Desired (`RD=0`) flag to manually traverse the internet's backbone, starting at the Root Name Servers (`198.41.0.4`) and chaining TLD and Authoritative referrals (via Glue Records) until the final IP is resolved.
*   **Manual Packet Serialization:** Implements custom binary packing and unpacking for DNS Headers, Questions, and Answer Records directly onto byte slices using Big-Endian formatting.
*   **DNS Pointer Decompression:** Implements a robust recursive parser to unwrap 16-bit compression pointers (detecting the `0xC0` MSB trap) to successfully reconstruct dynamically compressed domain names while preventing infinite loop vulnerabilities.
*   **Recursive Fallback:** Supports offloading the resolution work to a public resolver (like Google's `8.8.8.8`) via a simple CLI flag for lightweight querying.

## Usage

Run the resolver directly from the terminal using the Go CLI.

**Iterative Resolution (Default):**
Walks the DNS tree starting from the Root Servers.
```bash
go run main.go -domain=github.com
```

**Recursive Resolution:**
Sets the `RD=1` flag and queries Google's DNS to handle the traversal.
```bash
go run main.go -domain=github.com -recursive
```

## Future Roadmap

While the core IPv4 (A Record) resolution is fully operational, there are several systems-level features planned for future implementation:

*   **CNAME Resolution:** Add parsing and follow-through logic for Canonical Name (Type 5) aliases.
*   **Robust Authority Resolution:** Currently, the iterative loop relies on IP addresses being provided in the `Additionals` (Glue) section. If an Authoritative server is returned by name only, the resolver needs to pause, recursively resolve that Nameserver's IP from scratch, and then resume the original query.
*   **Local Caching:** Implement a local, TTL-respecting memory cache to store resolved IPs and TLD Nameserver IPs. This will significantly speed up subsequent queries and eliminate redundant network load on the Root Servers.
