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

func (r *Record) EncodeRecord() []byte {
	res := make([]byte, 0)

	res = append(res, r.Name...)

	res = binary.BigEndian.AppendUint16(res, r.Type)
	res = binary.BigEndian.AppendUint16(res, r.Class)

	res = binary.BigEndian.AppendUint32(res, r.TTL)

	res = binary.BigEndian.AppendUint16(res, r.DataLen)

	res = append(res, r.Data...)

	return res
}

func (h *Header) EncodeHeader() []byte {
	res := make([]byte, 12)
	binary.BigEndian.PutUint16(res[0:2], h.ID)
	binary.BigEndian.PutUint16(res[2:4], h.Flags)
	binary.BigEndian.PutUint16(res[4:6], h.QDCount)
	binary.BigEndian.PutUint16(res[6:8], h.ANCount)
	binary.BigEndian.PutUint16(res[8:10], h.NSCount)
	binary.BigEndian.PutUint16(res[10:12], h.ARCount)
	return res
}

func BuildQuery(msg DNSMessage) []byte {
	res := make([]byte, 0)
	res = append(res, msg.H.EncodeHeader()...)
	res = append(res, msg.Q.EncodeQuestion()...)
	return res
}

func NewSimpleQuery(id uint16, domain string) []byte {
	msg := DNSMessage{
		H: Header{
			ID:      id,
			QDCount: 1,
		},
		Q: Question{
			Name:  EncodeDomainName(domain),
			Type:  1, // Type A
			Class: 1, // Class IN
		},
	}

	// QR=0, OPCODE=0, AA=0, TC=0, RD=1, RA=0, Z=0, RCODE=0
	msg.H.SetFlags(0, 0, 0, 0, 1, 0, 0, 0)

	return BuildQuery(msg)
}
