package resolver

import "encoding/binary"

type DNSMessage struct {
	H Header
  Q Question
}

//     The Header Section
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type Header [12]byte

func (h *Header) SetID(id uint16) {
	binary.BigEndian.PutUint16(h[0:2], id)
}

func (h *Header) SetFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint8) {
	h[2] = (qr << 7) | (opcode << 3) | (aa << 2) | (tc << 1) | rd

	h[3] = (ra << 7) | (z << 4) | rcode
}

func (h *Header) SetQDCount(qd uint16) {
	binary.BigEndian.PutUint16(h[4:6], qd)
}

func (h *Header) SetANCount(an uint16) {
	binary.BigEndian.PutUint16(h[6:8], an)
}

func (h *Header) SetNSCount(ns uint16) {
	binary.BigEndian.PutUint16(h[8:10], ns)
}

func (h *Header) SetARCount(ar uint16) {
	binary.BigEndian.PutUint16(h[10:12], ar)
}

func (h *Header) ID() uint16 {
	return binary.BigEndian.Uint16(h[0:2])
}

func (h *Header) QDCount() uint16 {
	return binary.BigEndian.Uint16(h[4:6])
}

func (h *Header) ANCount() uint16 {
	return binary.BigEndian.Uint16(h[6:8])
}

func (h *Header) NSCount() uint16 {
	return binary.BigEndian.Uint16(h[8:10])
}

func (h *Header) ARCount() uint16 {
	return binary.BigEndian.Uint16(h[10:12])
}

func (h *Header) Flags() (qr, opcode, aa, tc, rd, ra, z, rcode uint8) {
	qr = (h[2] >> 7) & 0x01
	opcode = (h[2] >> 3) & 0x0F
	aa = (h[2] >> 2) & 0x01
	tc = (h[2] >> 1) & 0x01
	rd = h[2] & 0x01

	ra = (h[3] >> 7) & 0x01
	z = (h[3] >> 4) & 0x07
	rcode = h[3] & 0x0F

	return
}

// QUESTION SECTION
//                                 1  1  1  1  1  1
//   0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                                               |
// /                     QNAME                     /
// /                                               /
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QTYPE                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QCLASS                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//

type Question struct {
	Name  []byte
	Type  uint16
	Class uint16
}
