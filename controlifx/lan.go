package controlifx

import (
	"encoding"
	"encoding/binary"
	"fmt"
	_time "time"
)

// The recommended maximum number of messages to be sent to any one device
// every second over LAN.
const MessageRate = 20

// The LAN protocol header is always 36 bytes long.
const LanHeaderSize = 36

type BinaryEncoder interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type LanMessage struct {
	header  LanHeader
	payload BinaryEncoder
}

func (o *LanMessage) Payload(payload BinaryEncoder) {
	o.payload = payload

	o.updateSize()
}

func (o *LanMessage) updateSize() {
	size := LanHeaderSize

	if o.payload != nil {
		b, _ := o.payload.MarshalBinary()

		size += len(b)
	}

	o.header.frame.Size = uint16(size)
}

func (o LanMessage) MarshalBinary() (data []byte, err error) {
	header, err := o.header.MarshalBinary()
	if err != nil {
		return
	}

	var payload []byte

	if o.payload != nil {
		payload, err = o.payload.MarshalBinary()
		if err != nil {
			return
		}
	}

	data = append(header, payload...)

	return
}

func (o *LanMessage) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

type LanHeader struct {
	frame 		   LanHeaderFrame
	frameAddress   LanHeaderFrameAddress
	protocolHeader LanHeaderProtocolHeader
}

func (o LanHeader) MarshalBinary() (data []byte, err error) {
	frame, err := o.frame.MarshalBinary()
	if err != nil {
		return
	}

	frameAddress, err := o.frameAddress.MarshalBinary()
	if err != nil {
		return
	}

	protocolHeader, err := o.protocolHeader.MarshalBinary()
	if err != nil {
		return
	}

	data = append(frame, append(frameAddress, protocolHeader...)...)

	return
}

func (o *LanHeader) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

type LanHeaderFrame struct {
	Size   uint16
	Tagged bool
	Source uint32
}

func (o LanHeaderFrame) MarshalBinary() (data []byte, _ error) {
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

func (o *LanHeaderFrame) UnmarshalBinary(data []byte) error {
	o.Size = binary.LittleEndian.Uint16(data[:2])
	o.Tagged = (data[2] >> 5) & 0x1 == 1
	o.Source = binary.LittleEndian.Uint32(data[4:])

	return nil
}

type LanHeaderFrameAddress struct {
	Target      uint64
	AckRequired bool
	ResRequired bool
	Sequence    uint8
}

func (o LanHeaderFrameAddress) MarshalBinary() (data []byte, _ error) {
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

func (o *LanHeaderFrameAddress) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

type LanHeaderProtocolHeader struct {
	Type uint16
}

func (o LanHeaderProtocolHeader) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 12)

	// Type.
	binary.LittleEndian.PutUint16(data[8:10], o.Type)

	return
}

type label struct {
	value string
}

func (o label) MarshalBinary() (data []byte, err error) {
	const Size = 32

	if len(o.value) > Size {
		err = fmt.Errorf("label '%s' has length %d > %d", o.value, len(o.value), Size)
		return
	}

	data = append([]byte(o.value), make([]byte, Size - len(o.value))...)

	return
}

type port struct {
	value uint32
}

func (o port) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 4)

	binary.LittleEndian.PutUint32(data, o.value)

	return
}

type powerLevel struct {
	value uint16
}

func (o powerLevel) MarshalBinary() (data []byte, err error) {
	if o.value != 0 && o.value != 65535 {
		err = fmt.Errorf("level %d is not 0 or 65535", o.value)
		return
	}

	data = make([]byte, 2)

	binary.LittleEndian.PutUint16(data, o.value)

	return
}

type time struct {
	value uint64
}

func (o time) Time() _time.Time {
	// Check if value is over the max int64 size.
	if o.value > 9223372036854775807 {
		return _time.Time{}
	}

	return _time.Unix(0, int64(o.value))
}

func (o time) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 8)

	binary.LittleEndian.PutUint64(data, o.value)

	return
}

type LanDeviceMessageBuilder struct {
	Source      uint32
	Target      uint64
	AckRequired bool
	ResRequired bool
	Sequence    uint8
}

func (o LanDeviceMessageBuilder) Tagged() bool {
	return o.Target > 0
}

func (o LanDeviceMessageBuilder) buildNormalMessageOfType(t uint16) LanMessage {
	return LanMessage{
		header:LanHeader{
			frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:o.Tagged(),
				Source:o.Source,
			},
			frameAddress:LanHeaderFrameAddress{
				Target:o.Target,
				AckRequired:o.AckRequired,
				ResRequired:o.ResRequired,
				Sequence:o.Sequence,
			},
			protocolHeader:LanHeaderProtocolHeader{
				Type:t,
			},
		},
	}
}

func (o LanDeviceMessageBuilder) GetService() LanMessage {
	const Type = 2

	return o.buildNormalMessageOfType(Type)
}

type StateServiceLanMessage struct {
	Service uint8
	Port    uint32
}

func (o LanDeviceMessageBuilder) GetHostInfo() LanMessage {
	const Type = 12

	return o.buildNormalMessageOfType(Type)
}

type StateHostInfoLanMessage struct {
	Signal float32
	Tx     uint32
	Rx     uint32
}

func (o LanDeviceMessageBuilder) GetHostFirmware() LanMessage {
	const Type = 14

	return o.buildNormalMessageOfType(Type)
}

type StateHostFirmwareLanMessage struct {
	Build   uint64
	Version uint32
}

func (o LanDeviceMessageBuilder) GetWifiInfo() LanMessage {
	const Type = 16

	return o.buildNormalMessageOfType(Type)
}

type StateWifiInfoLanMessage struct {
	Signal float32
	Tx     uint32
	Rx     uint32
}

func (o LanDeviceMessageBuilder) GetWifiFirmware() LanMessage {
	const Type = 18

	return o.buildNormalMessageOfType(Type)
}

type StateWifiFirmwareLanMessage struct {
	Build   uint64
	Version uint32
}

func (o LanDeviceMessageBuilder) GetPower() LanMessage {
	const Type = 20

	return o.buildNormalMessageOfType(Type)
}

type SetPowerLanMessage struct {
	Level powerLevel
}

func (o SetPowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}

func (o *SetPowerLanMessage) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

func (o LanDeviceMessageBuilder) SetPower(payload *SetPowerLanMessage) LanMessage {
	const Type = 21

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type StatePowerLanMessage struct {
	Level powerLevel
}

func (o LanDeviceMessageBuilder) GetLabel() LanMessage {
	const Type = 23

	return o.buildNormalMessageOfType(Type)
}

type SetLabelLanMessage struct {
	label label
}

func (o SetLabelLanMessage) MarshalBinary() ([]byte, error) {
	return o.label.MarshalBinary()
}

func (o *SetLabelLanMessage) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

func (o LanDeviceMessageBuilder) SetLabel(payload *SetLabelLanMessage) LanMessage {
	const Type = 24

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type StateLabelLanMessage struct {
	label label
}

func (o LanDeviceMessageBuilder) GetVersion() LanMessage {
	const Type = 32

	return o.buildNormalMessageOfType(Type)
}

type StateVersionLanMessage struct {
	vendor  uint32
	product uint32
	version uint32
}

func (o LanDeviceMessageBuilder) GetInfo() LanMessage {
	const Type = 34

	return o.buildNormalMessageOfType(Type)
}

type StateInfoLanMessage struct {
	time     time
	uptime   uint64
	downtime uint64
}

func (o LanDeviceMessageBuilder) GetLocation() LanMessage {
	const Type = 48

	return o.buildNormalMessageOfType(Type)
}

type StateLocationLanMessage struct {
	location  [16]byte
	label     label
	updatedAt time
}

func (o LanDeviceMessageBuilder) GetGroup() LanMessage {
	const Type = 51

	return o.buildNormalMessageOfType(Type)
}

type StateGroupLanMessage struct {
	group     [16]byte
	label     label
	updatedAt time
}

type EchoRequestLanMessage struct {
	payload [64]byte
}

func (o EchoRequestLanMessage) MarshalBinary() ([]byte, error) {
	return o.payload[:], nil
}

func (o *EchoRequestLanMessage) UnmarshalBinary(data []byte) error {
	// TODO: implement

	return nil
}

func (o LanDeviceMessageBuilder) EchoRequest(payload *EchoRequestLanMessage) LanMessage {
	const Type = 58

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type EchoResponseLanMessage struct {
	payload [64]byte
}
