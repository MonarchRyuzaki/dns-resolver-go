package resolver

import "encoding/binary"

func DecodeHeader(data []byte) Header {
	return Header{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		QDCount: binary.BigEndian.Uint16(data[4:6]),
		ANCount: binary.BigEndian.Uint16(data[6:8]),
		NSCount: binary.BigEndian.Uint16(data[8:10]),
		ARCount: binary.BigEndian.Uint16(data[10:12]),
	}
}

func DecodeQuestion(data []byte, offset int) (Question, int) {
	name, newOffset := DecodeDomainName(data, offset)

	q := Question{
		Name:  name,
		Type:  binary.BigEndian.Uint16(data[newOffset : newOffset+2]),
		Class: binary.BigEndian.Uint16(data[newOffset+2 : newOffset+4]),
	}

	return q, newOffset + 4
}

// DecodeResponse takes the raw UDP bytes and parses them into a full DNSMessage.
func DecodeResponse(data []byte) DNSMessage {
	msg := DNSMessage{}

	msg.H = DecodeHeader(data[:12])

	offset := 12
	msg.Q, offset = DecodeQuestion(data, offset)

	// 3. Now 'offset' is pointing directly at the Answer section.
	// Parse the Record Name (Watch out for DNS Compression pointers!)
	// TODO: decode name, increment offset

	// Parse Type, Class, TTL, and DataLen (10 bytes total)
	// TODO: Read uint16 Type, uint16 Class, uint32 TTL, uint16 DataLen

	// Parse Data
	// TODO: Read the next 'DataLen' bytes as the actual IP address bytes

	// TODO: assign the parsed Record to msg.Answer

	return msg
}

// DecodeDomainName is a helper you'll absolutely need to handle DNS Compression.
// It returns the decoded string and the new offset.
func DecodeDomainName(data []byte, offset int) ([]byte, int) {
	cnt := 100
	res := make([]byte, 0)
	tempOffset := offset

	jumped := false
	returnOffset := offset

	for cnt > 0 {
		readByte := data[tempOffset]

		if readByte&0xC0 == 0xC0 {
			if !jumped {
				returnOffset = tempOffset + 2
				jumped = true
			}
			pointerBytes := binary.BigEndian.Uint16(data[tempOffset : tempOffset+2])
			jumpOffset := pointerBytes & 0x3FFF

			tempOffset = int(jumpOffset)
			cnt--
			continue
		}

		if readByte == 0x00 {
			if !jumped {
				returnOffset = tempOffset + 1
			}
			break
		}

		length := int(readByte)

		res = append(res, byte(length))
		res = append(res, data[tempOffset+1:tempOffset+1+length]...)

		tempOffset += 1 + length
		cnt--
	}

	res = append(res, 0x00)

	return res, returnOffset
}
