package resolver

// DecodeResponse takes the raw UDP bytes and parses out the first Answer Record.
func DecodeResponse(data []byte) *Record {
	// 1. Unpack the first 12 bytes into your Header type
	// TODO: Create a Header, copy data[0:12] into it. Check ANCount to ensure you have answers!

	// 2. We need to skip the Question section.
	// The Question starts at byte 12. You must read the variable-length QNAME
	// until you hit a 0x00 byte, then jump over the 4 bytes for QTYPE and QCLASS.
	offset := 12
	_ = offset
	// TODO: increment 'offset' past the Question section

	// 3. Now 'offset' is pointing directly at the Answer section.
	// Parse the Record Name (Watch out for DNS Compression pointers!)
	// TODO: decode name, increment offset

	// Parse Type, Class, TTL, and DataLen (10 bytes total)
	// TODO: Read uint16 Type, uint16 Class, uint32 TTL, uint16 DataLen

	// Parse Data
	// TODO: Read the next 'DataLen' bytes as the actual IP address bytes

	return &Record{}
}

// DecodeDomainName is a helper you'll absolutely need to handle DNS Compression.
// It returns the decoded string and the new offset.
func DecodeDomainName(data []byte, offset int) ([]byte, int) {
	// TODO: Implement parsing
	// Trap: If the top two bits of the length byte are 11 (e.g. byte >= 192 or 0xC0),
	// it's a pointer to somewhere else in the packet!
	return nil, offset
}
