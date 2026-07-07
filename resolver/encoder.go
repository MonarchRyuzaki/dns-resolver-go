// Package resolver provides the core types and logic for manually
// crafting and parsing DNS UDP packets according to RFC 1035
package resolver

import (
	"encoding/binary"
	"strings"
)

func EncodeDomainName(domain string) []byte {
	splittedString := strings.Split(domain, ".")
	res := make([]byte, 0)
	for _, v := range splittedString {
		res = append(res, byte(len(v)))

		res = append(res, []byte(v)...)
	}

	res = append(res, 0x00)

	return res
}

func (Q *Question) EncodeQuestion() []byte {
	res := make([]byte, 0)
	res = append(res, Q.Name...)
	res = binary.BigEndian.AppendUint16(res, Q.Type)
	res = binary.BigEndian.AppendUint16(res, Q.Class)

	return res
}

func BuildQuery(msg DNSMessage) []byte {
	res := make([]byte, 0)
	res = append(res, msg.H[:]...)
	res = append(res, msg.Q.EncodeQuestion()...)
	return res
}

func NewSimpleQuery(id uint16, domain string) []byte {
	msg := DNSMessage{
		H: Header{},
		Q: Question{
			Name:  EncodeDomainName(domain),
			Type:  1, // Type A
			Class: 1, // Class IN
		},
	}

	msg.H.SetID(id)
	// QR=0, OPCODE=0, AA=0, TC=0, RD=1, RA=0, Z=0, RCODE=0
	msg.H.SetFlags(0, 0, 0, 0, 1, 0, 0, 0)
	msg.H.SetQDCount(1)

	return BuildQuery(msg)
}
