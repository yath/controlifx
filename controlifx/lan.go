package controlifx

import "encoding/binary"

// The recommended maximum number of messages to be sent to any one device
// every second.
const MessageRate = 20

type LanMessage struct {
	header  LanHeader
	payload LanPayload
}

type LanHeader struct {
	frame 		   LanHeaderFrame
	frameAddress   LanHeaderFrameAddress
	protocolHeader LanHeaderProtocolHeader
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
	Target      uint64
	AckRequired bool
	ResRequired bool
	Sequence    uint8
}

func (o *LanHeaderFrameAddress) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 16)

	littleEndianPutUint48 := func(b []byte, v uint64) {
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
	}

	// Target.
	if o.Target <= 0xffffff {
		littleEndianPutUint48(data[:6], o.Target)
	} else {
		binary.LittleEndian.PutUint64(data[:8], o.Target)
	}

	// 0000 0000

	// AckRequired.
	if o.AckRequired {
		data[14] |= 0x02
	}
	// 0000 00?0

	// ResRequired.
	if o.ResRequired {
		data[14] |= 0x01
	}
	// 0000 00??

	// Sequence.
	data[15] = o.Sequence

	return
}

type LanHeaderProtocolHeader struct {
	Type uint16
}

func (o *LanHeaderProtocolHeader) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Type.
	binary.LittleEndian.PutUint16(data[8:10], o.Type)

	return
}

type LanPayload struct {

}
