package controlifx

import "encoding/binary"

type LanMessage struct {

}

type LanHeader struct {

}

type LanHeaderFrame struct {
	Size   uint16
	Tagged bool
	Source uint32
}

func (o *LanHeaderFrame) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 8)

	// Size.
	binary.LittleEndian.PutUint16(data[:2], o.Size)

	// Origin.
	// 0000 0000  0000 0000

	// Tagged.
	if o.Tagged {
		data[2] |= 0x20
	}
	// 00?0 0000  0000 0000

	// Addressable (1) | Protocol (1024).
	data[2] |= 0x18
	// 00?1 1000  0000 0000

	// Source.
	binary.LittleEndian.PutUint32(data[4:], o.Source)

	return
}

type LanHeaderFrameAddress struct {

}

type LanHeaderProtocolHeader struct {

}

type LanHeaderPayload struct {

}
