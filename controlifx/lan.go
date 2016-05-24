package controlifx

import (
	"encoding"
	"encoding/binary"
	"fmt"
	_time "time"
	"math"
)

// The recommended maximum number of messages to be sent to any one device
// every second over LAN.
const MessageRate = 20

// The LAN protocol header is always 36 bytes long.
const LanHeaderSize = 36

type SendableLanMessage struct {
	header  LanHeader
	payload encoding.BinaryMarshaler
}

func (o *SendableLanMessage) Payload(payload encoding.BinaryMarshaler) {
	o.payload = payload

	o.updateSize()
}

func (o *SendableLanMessage) updateSize() {
	size := LanHeaderSize

	if o.payload != nil {
		b, _ := o.payload.MarshalBinary()

		size += len(b)
	}

	o.header.frame.Size = uint16(size)
}

func (o SendableLanMessage) MarshalBinary() (data []byte, err error) {
	// Header.
	header, err := o.header.MarshalBinary()
	if err != nil {
		return
	}

	// Payload.
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

type ReceivableLanMessage struct {
	header  LanHeader
	payload encoding.BinaryUnmarshaler
}

func (o *ReceivableLanMessage) UnmarshalBinary(data []byte) error {
	// Header.
	o.header = LanHeader{}
	if err := o.header.UnmarshalBinary(data[:LanHeaderSize]); err != nil {
		return err
	}

	// Payload.
	o.payload = NewPayloadOfType(o.header.protocolHeader.Type)
	return o.payload.UnmarshalBinary(data[LanHeaderSize:])
}

type LanHeader struct {
	frame 		   LanHeaderFrame
	frameAddress   LanHeaderFrameAddress
	protocolHeader LanHeaderProtocolHeader
}

func (o LanHeader) MarshalBinary() (data []byte, err error) {
	// Frame.
	frame, err := o.frame.MarshalBinary()
	if err != nil {
		return
	}

	// Frame address.
	frameAddress, err := o.frameAddress.MarshalBinary()
	if err != nil {
		return
	}

	// Protocol header.
	protocolHeader, err := o.protocolHeader.MarshalBinary()
	if err != nil {
		return
	}

	data = append(frame, append(frameAddress, protocolHeader...)...)

	return
}

func (o *LanHeader) UnmarshalBinary(data []byte) error {
	// Frame.
	o.frame = LanHeaderFrame{}
	if err := o.frame.UnmarshalBinary(data[:8]); err != nil {
		return err
	}

	// Frame address.
	o.frameAddress = LanHeaderFrameAddress{}
	if err := o.frameAddress.UnmarshalBinary(data[8:24]); err != nil {
		return err
	}

	// Protocol header.
	o.protocolHeader = LanHeaderProtocolHeader{}
	return o.protocolHeader.UnmarshalBinary(data[24:])
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
	// Size.
	o.Size = binary.LittleEndian.Uint16(data[:2])

	// Tagged.
	o.Tagged = (data[2] >> 5) & 0x1 == 1

	// Source.
	o.Source = binary.LittleEndian.Uint32(data[4:])

	return nil
}

type LanHeaderFrameAddress struct {
	Target      uint64
	// Not part of the protocol, but specifies if the Target is a 6 byte MAC
	// address for binary marshaling purposes.
	TargetMac   bool
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
	if o.TargetMac {
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
	littleEndianUint48 := func(b []byte) uint64 {
		return uint64(b[0]) | uint64(b[1]) << 8 | uint64(b[2]) << 16 |
				uint64(b[3]) << 24 | uint64(b[4]) << 32 | uint64(b[5]) << 40
	}

	// Target. There may be issues if an 8 byte device address is used where the
	// last 2 bytes are 0x00, but this seems pretty rare, and the Target from a
	// response is probably not that important anyway. Really, there is no way
	// for the client to know if it's a MAC address or some other 8 byte long
	// address unless it looks at the sequence to map it with the original
	// request.
	if binary.LittleEndian.Uint16(data[6:8]) == 0 {
		o.Target = binary.LittleEndian.Uint64(data[:8])
	} else {
		o.Target = littleEndianUint48(data[:6])
	}

	// 0000 00??

	// AckRequired.
	o.AckRequired = ((data[14] >> 1) & 0x01 == 1)

	// ResRequired.
	o.ResRequired = (data[14] & 0x01 == 1)

	// Sequence.
	o.Sequence = byte(data[15])

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

func (o *LanHeaderProtocolHeader) UnmarshalBinary(data []byte) error {
	// Type.
	o.Type = binary.LittleEndian.Uint16(data[8:10])

	return nil
}

type label string

func (o label) MarshalBinary() (data []byte, err error) {
	const Size = 32

	if len(o) > Size {
		err = fmt.Errorf("label '%s' has length %d > %d", o, len(o), Size)
		return
	}

	data = append([]byte(o), make([]byte, Size - len(o))...)

	return
}

type port uint32

func (o port) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 4)

	binary.LittleEndian.PutUint32(data, uint32(o))

	return
}

type powerLevel uint16

func (o powerLevel) MarshalBinary() (data []byte, err error) {
	if o != 0 && o != 65535 {
		err = fmt.Errorf("level %d is not 0 or 65535", o)
		return
	}

	data = make([]byte, 2)

	binary.LittleEndian.PutUint16(data, uint16(o))

	return
}

type time uint64

func (o time) Time() _time.Time {
	// Check if value is over the max int64 size.
	if o > 9223372036854775807 {
		return _time.Time{}
	}

	return _time.Unix(0, int64(o))
}

func (o time) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(o))

	return
}

func NewPayloadOfType(t uint16) encoding.BinaryUnmarshaler {
	const (
		StateService      = 3
		StateHostInfo     = 13
		StateHostFirmware = 15
		StateWifiInfo     = 17
		StateWifiFirmware = 19
		StatePower        = 22
		StateLabel        = 25
		StateVersion      = 33
		StateInfo         = 35
		StateLocation     = 50
		StateGroup        = 53
		EchoResponse      = 59
	)

	switch t {
	case StateService:
		return &StateServiceLanMessage{}
	case StateHostInfo:
		return &StateHostInfoLanMessage{}
	case StateHostFirmware:
		return &StateHostFirmwareLanMessage{}
	case StateWifiInfo:
		return &StateWifiInfoLanMessage{}
	case StateWifiFirmware:
		return &StateWifiFirmwareLanMessage{}
	case StatePower:
		return &StatePowerLanMessage{}
	case StateLabel:
		return &StateLabelLanMessage{}
	case StateVersion:
		return &StateVersionLanMessage{}
	case StateInfo:
		return &StateInfoLanMessage{}
	case StateLocation:
		return &StateLocationLanMessage{}
	case StateGroup:
		return &StateGroupLanMessage{}
	case EchoResponse:
		return &EchoResponseLanMessage{}
	default:
		return nil
	}
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

func (o LanDeviceMessageBuilder) buildNormalMessageOfType(t uint16) SendableLanMessage {
	return SendableLanMessage{
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

func (o LanDeviceMessageBuilder) GetService() SendableLanMessage {
	const Type = 2

	return o.buildNormalMessageOfType(Type)
}

type StateServiceLanMessage struct {
	Service uint8
	Port    uint32
}

func (o *StateServiceLanMessage) UnmarshalBinary(data []byte) error {
	// Service.
	o.Service = uint8(data[0])

	// Port.
	o.Port = binary.LittleEndian.Uint32(data[1:])

	return nil
}

func (o LanDeviceMessageBuilder) GetHostInfo() SendableLanMessage {
	const Type = 12

	return o.buildNormalMessageOfType(Type)
}

type StateHostInfoLanMessage struct {
	Signal float32
	Tx     uint32
	Rx     uint32
}

func (o *StateHostInfoLanMessage) UnmarshalBinary(data []byte) error {
	// Signal.
	o.Signal = math.Float32frombits(binary.LittleEndian.Uint32(data[:4]))

	// Tx.
	o.Tx = binary.LittleEndian.Uint32(data[4:8])

	// Rx.
	o.Rx = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func (o LanDeviceMessageBuilder) GetHostFirmware() SendableLanMessage {
	const Type = 14

	return o.buildNormalMessageOfType(Type)
}

type StateHostFirmwareLanMessage struct {
	Build   uint64
	Version uint32
}

func (o *StateHostFirmwareLanMessage) UnmarshalBinary(data []byte) error {
	// Build.
	o.Build = binary.LittleEndian.Uint64(data[:8])

	// Version.
	o.Version = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func (o LanDeviceMessageBuilder) GetWifiInfo() SendableLanMessage {
	const Type = 16

	return o.buildNormalMessageOfType(Type)
}

type StateWifiInfoLanMessage struct {
	Signal float32
	Tx     uint32
	Rx     uint32
}

func (o *StateWifiInfoLanMessage) UnmarshalBinary(data []byte) error {
	// Signal.
	o.Signal = math.Float32frombits(binary.LittleEndian.Uint32(data[:4]))

	// Tx.
	o.Tx = binary.LittleEndian.Uint32(data[4:8])

	// Rx.
	o.Rx = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func (o LanDeviceMessageBuilder) GetWifiFirmware() SendableLanMessage {
	const Type = 18

	return o.buildNormalMessageOfType(Type)
}

type StateWifiFirmwareLanMessage struct {
	Build   uint64
	Version uint32
}

func (o *StateWifiFirmwareLanMessage) UnmarshalBinary(data []byte) error {
	// Build.
	o.Build = binary.LittleEndian.Uint64(data[:8])

	// Version.
	o.Version = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func (o LanDeviceMessageBuilder) GetPower() SendableLanMessage {
	const Type = 20

	return o.buildNormalMessageOfType(Type)
}

type SetPowerLanMessage struct {
	Level powerLevel
}

func (o SetPowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}

func (o LanDeviceMessageBuilder) SetPower(payload *SetPowerLanMessage) SendableLanMessage {
	const Type = 21

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type StatePowerLanMessage struct {
	Level powerLevel
}

func (o *StatePowerLanMessage) UnmarshalBinary(data []byte) error {
	o.Level = powerLevel(binary.LittleEndian.Uint16(data))

	return nil
}

func (o LanDeviceMessageBuilder) GetLabel() SendableLanMessage {
	const Type = 23

	return o.buildNormalMessageOfType(Type)
}

type SetLabelLanMessage struct {
	label label
}

func (o SetLabelLanMessage) MarshalBinary() ([]byte, error) {
	return o.label.MarshalBinary()
}

func (o LanDeviceMessageBuilder) SetLabel(payload *SetLabelLanMessage) SendableLanMessage {
	const Type = 24

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type StateLabelLanMessage struct {
	label label
}

func (o *StateLabelLanMessage) UnmarshalBinary(data []byte) error {
	o.label = label(data)

	return nil
}

func (o LanDeviceMessageBuilder) GetVersion() SendableLanMessage {
	const Type = 32

	return o.buildNormalMessageOfType(Type)
}

type StateVersionLanMessage struct {
	vendor  uint32
	product uint32
	version uint32
}

func (o *StateVersionLanMessage) UnmarshalBinary(data []byte) error {
	// Vendor.
	o.vendor = binary.LittleEndian.Uint32(data[:4])

	// Product.
	o.product = binary.LittleEndian.Uint32(data[4:12])

	// Version.
	o.version = binary.LittleEndian.Uint32(data[12:])

	return nil
}

func (o LanDeviceMessageBuilder) GetInfo() SendableLanMessage {
	const Type = 34

	return o.buildNormalMessageOfType(Type)
}

type StateInfoLanMessage struct {
	time     time
	uptime   uint64
	downtime uint64
}

func (o *StateInfoLanMessage) UnmarshalBinary(data []byte) error {
	// Time.
	o.time = time(binary.LittleEndian.Uint64(data[:8]))

	// Uptime.
	o.uptime = binary.LittleEndian.Uint64(data[8:24])

	// Downtime.
	o.downtime = binary.LittleEndian.Uint64(data[24:])

	return nil
}

func (o LanDeviceMessageBuilder) GetLocation() SendableLanMessage {
	const Type = 48

	return o.buildNormalMessageOfType(Type)
}

type StateLocationLanMessage struct {
	location  [16]byte
	label     label
	updatedAt time
}

func (o *StateLocationLanMessage) UnmarshalBinary(data []byte) error {
	// Location.
	copy(o.location[:], data[:16])

	// Label.
	o.label = label(data[16:48])

	// Updated at.
	o.updatedAt = time(binary.LittleEndian.Uint64(data[48:]))

	return nil
}

func (o LanDeviceMessageBuilder) GetGroup() SendableLanMessage {
	const Type = 51

	return o.buildNormalMessageOfType(Type)
}

type StateGroupLanMessage struct {
	group     [16]byte
	label     label
	updatedAt time
}

func (o *StateGroupLanMessage) UnmarshalBinary(data []byte) error {
	// Group.
	copy(o.group[:], data[:16])

	// Label.
	o.label = label(data[16:48])

	// Updated at.
	o.updatedAt = time(binary.LittleEndian.Uint64(data[48:]))

	return nil
}

type EchoRequestLanMessage struct {
	payload [64]byte
}

func (o EchoRequestLanMessage) MarshalBinary() ([]byte, error) {
	return o.payload[:], nil
}

func (o LanDeviceMessageBuilder) EchoRequest(payload *EchoRequestLanMessage) SendableLanMessage {
	const Type = 58

	msg := o.buildNormalMessageOfType(Type)

	msg.Payload(payload)

	return msg
}

type EchoResponseLanMessage struct {
	payload [64]byte
}

func (o *EchoResponseLanMessage) UnmarshalBinary(data []byte) error {
	copy(o.payload[:], data[:64])

	return nil
}
