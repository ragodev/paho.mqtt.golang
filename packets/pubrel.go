package packets

import (
	"bytes"
	"net"
)

// Pubrel is the Variable Header definition for a Pubrel control packet
type Pubrel struct {
	PacketID   uint16
	ReasonCode byte
	IDVP       IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubrel) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	success := r.Len() == 2
	p.PacketID, err = readUint16(r)
	if err != nil {
		return 0, err
	}
	if !success {
		p.ReasonCode, err = r.ReadByte()
		if err != nil {
			return 0, err
		}

		idvpLen, err := p.IDVP.Unpack(r, PUBACK)
		if err != nil {
			return 0, err
		}

		return idvpLen + 3, nil
	}
	return 2, nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Pubrel) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(p.PacketID, &b)
	b.WriteByte(p.ReasonCode)
	idvp := p.IDVP.Pack(PUBREL)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp}
}
