package controlifx

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
	"math"
)

// The maximum recommended number of messages to be sent to any one device every
// second.
const MessageRate = 20

// The LAN protocol header is always 36 bytes long.
const LanHeaderSize = 36

type SendableLanMessage struct {
	Header  LanHeader
	Payload encoding.BinaryMarshaler
}

func (o *SendableLanMessage) updateSize() {
	size := LanHeaderSize

	if o.Payload != nil {
		b, _ := o.Payload.MarshalBinary()

		size += len(b)
	}

	o.Header.Frame.Size = uint16(size)
}

func (o SendableLanMessage) MarshalBinary() (data []byte, err error) {
	// Header.
	header, err := o.Header.MarshalBinary()
	if err != nil {
		return
	}

	// Payload.
	var payload []byte

	if o.Payload != nil {
		payload, err = o.Payload.MarshalBinary()
		if err != nil {
			return
		}
	}

	data = append(header, payload...)

	return
}

type ReceivableLanMessage struct {
	Header  LanHeader
	Payload encoding.BinaryUnmarshaler
}

func (o *ReceivableLanMessage) UnmarshalBinary(data []byte) error {
	// Header.
	o.Header = LanHeader{}
	if err := o.Header.UnmarshalBinary(data[:LanHeaderSize]); err != nil {
		return err
	}

	// Payload.
	payload, err := NewReceivablePayloadOfType(o.Header.ProtocolHeader.Type)
	if err != nil {
		return err
	}

	o.Payload = payload

	return o.Payload.UnmarshalBinary(data[LanHeaderSize:])
}

type LanHeader struct {
	Frame          LanHeaderFrame
	FrameAddress   LanHeaderFrameAddress
	ProtocolHeader LanHeaderProtocolHeader
}

func (o LanHeader) MarshalBinary() (data []byte, err error) {
	// Frame.
	frame, err := o.Frame.MarshalBinary()
	if err != nil {
		return
	}

	// Frame address.
	frameAddress, err := o.FrameAddress.MarshalBinary()
	if err != nil {
		return
	}

	// Protocol header.
	protocolHeader, err := o.ProtocolHeader.MarshalBinary()
	if err != nil {
		return
	}

	data = append(frame, append(frameAddress, protocolHeader...)...)

	return
}

func (o *LanHeader) UnmarshalBinary(data []byte) error {
	// Frame.
	o.Frame = LanHeaderFrame{}
	if err := o.Frame.UnmarshalBinary(data[:8]); err != nil {
		return err
	}

	// Frame address.
	o.FrameAddress = LanHeaderFrameAddress{}
	if err := o.FrameAddress.UnmarshalBinary(data[8:24]); err != nil {
		return err
	}

	// Protocol header.
	o.ProtocolHeader = LanHeaderProtocolHeader{}
	return o.ProtocolHeader.UnmarshalBinary(data[24:])
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
		data[3] |= 0x20
	}
	// 00?0 0000  0000 0000

	// Addressable (1) | Protocol (1024).
	data[3] |= 0x14
	// 0000 0000  00?1 0100

	// Source.
	binary.LittleEndian.PutUint32(data[4:], o.Source)

	return
}

func (o *LanHeaderFrame) UnmarshalBinary(data []byte) error {
	// Size.
	o.Size = binary.LittleEndian.Uint16(data[:2])

	// Tagged.
	o.Tagged = (data[3] >> 5) & 0x1 == 1

	// Source.
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

	// Big endian.
	putUint48 := func(b []byte, v uint64) {
		b[0] = byte(v >> 40)
		b[1] = byte(v >> 32)
		b[2] = byte(v >> 24)
		b[3] = byte(v >> 16)
		b[4] = byte(v >> 8)
		b[5] = byte(v)
	}

	// Target.
	putUint48(data[:6], o.Target)

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
	// Big endian.
	uint48 := func(b []byte) uint64 {
		return uint64(b[5]) | uint64(b[4])<<8 | uint64(b[3])<<16 | uint64(b[2])<<24 |
			uint64(b[1])<<32 | uint64(b[0])<<40
	}
	o.Target = uint48(data[:6])

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

type Label string

func (o Label) MarshalBinary() (data []byte, err error) {
	const Size = 32

	if len(o) > Size {
		err = fmt.Errorf("label '%s' has length %d > %d", o, len(o), Size)
		return
	}

	data = append([]byte(o), make([]byte, Size - len(o))...)

	return
}

type Port uint32

func (o Port) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 4)

	binary.LittleEndian.PutUint32(data, uint32(o))

	return
}

type PowerLevel uint16

func (o PowerLevel) MarshalBinary() (data []byte, err error) {
	if o != 0 && o != 65535 {
		err = fmt.Errorf("level %d is not 0 or 65535", o)
		return
	}

	data = make([]byte, 2)

	binary.LittleEndian.PutUint16(data[:], uint16(o))

	return
}

type Time uint64

func (o Time) Time() (time.Time, error) {
	// Check if value is over the max int64 size.
	if o > math.MaxInt64 {
		return time.Time{}, fmt.Errorf("time %d exceeds int64 max value", o)
	}

	return time.Unix(0, int64(o)), nil
}

func (o Time) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(o))

	return
}

func NewReceivablePayloadOfType(t uint16) (encoding.BinaryUnmarshaler, error) {
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
		Acknowledgement   = 45
		StateLocation     = 50
		StateGroup        = 53
		EchoResponse      = 59
		LightState        = 107
		LightStatePower   = 118
	)

	switch t {
	case StateService:
		return &StateServiceLanMessage{}, nil
	case StateHostInfo:
		return &StateHostInfoLanMessage{}, nil
	case StateHostFirmware:
		return &StateHostFirmwareLanMessage{}, nil
	case StateWifiInfo:
		return &StateWifiInfoLanMessage{}, nil
	case StateWifiFirmware:
		return &StateWifiFirmwareLanMessage{}, nil
	case StatePower:
		return &StatePowerLanMessage{}, nil
	case StateLabel:
		return &StateLabelLanMessage{}, nil
	case StateVersion:
		return &StateVersionLanMessage{}, nil
	case StateInfo:
		return &StateInfoLanMessage{}, nil
	case Acknowledgement:
		return &AcknowledgementLanMessage{}, nil
	case StateLocation:
		return &StateLocationLanMessage{}, nil
	case StateGroup:
		return &StateGroupLanMessage{}, nil
	case EchoResponse:
		return &EchoResponseLanMessage{}, nil
	case LightState:
		return &LightStatePowerLanMessage{}, nil
	case LightStatePower:
		return &LightStatePowerLanMessage{}, nil
	default:
		return nil, errors.New("cannot create new payload of type; is it binary encodable?")
	}
}

type LanDeviceMessageBuilder struct {
	source      uint32
	target      uint64
	AckRequired bool
	ResRequired bool
	Sequence    uint8
}

func (o LanDeviceMessageBuilder) Tagged() bool {
	return o.target > 0
}

func (o LanDeviceMessageBuilder) buildNormalMessageOfType(t uint16) SendableLanMessage {
	return SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:o.Tagged(),
				Source:o.source,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:o.target,
				AckRequired:o.AckRequired,
				ResRequired:o.ResRequired,
				Sequence:o.Sequence,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:t,
			},
		},
	}
}

const (
	GetServiceType      = 2
	GetHostInfoType     = 12
	GetHostFirmwareType = 14
	GetWifiInfoType     = 16
	GetWifiFirmwareType = 18
	GetPowerType        = 20
	SetPowerType        = 21
	GetLabelType        = 23
	SetLabelType        = 24
	GetVersionType      = 32
	GetInfoType         = 34
	GetLocationType     = 48
	GetGroupType        = 51
	EchoRequestType     = 58
	LightGetType        = 101
	LightSetColorType   = 102
	LightGetPowerType   = 116
	LightSetPowerType   = 117
)

func (o LanDeviceMessageBuilder) GetService() SendableLanMessage {
	return o.buildNormalMessageOfType(GetServiceType)
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
	return o.buildNormalMessageOfType(GetHostInfoType)
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
	return o.buildNormalMessageOfType(GetHostFirmwareType)
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
	return o.buildNormalMessageOfType(GetWifiInfoType)
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
	return o.buildNormalMessageOfType(GetWifiFirmwareType)
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
	return o.buildNormalMessageOfType(GetPowerType)
}

type SetPowerLanMessage struct {
	Level PowerLevel
}

func (o SetPowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}

func (o LanDeviceMessageBuilder) SetPower(payload SetPowerLanMessage) SendableLanMessage {
	msg := o.buildNormalMessageOfType(SetPowerType)

	msg.Payload = payload
	msg.updateSize()

	return msg
}

type StatePowerLanMessage struct {
	Level PowerLevel
}

func (o *StatePowerLanMessage) UnmarshalBinary(data []byte) error {
	o.Level = PowerLevel(binary.LittleEndian.Uint16(data))

	return nil
}

func (o LanDeviceMessageBuilder) GetLabel() SendableLanMessage {
	return o.buildNormalMessageOfType(GetLabelType)
}

type SetLabelLanMessage struct {
	Label Label
}

func (o SetLabelLanMessage) MarshalBinary() ([]byte, error) {
	return o.Label.MarshalBinary()
}

func (o LanDeviceMessageBuilder) SetLabel(payload SetLabelLanMessage) SendableLanMessage {
	msg := o.buildNormalMessageOfType(SetLabelType)

	msg.Payload = payload
	msg.updateSize()

	return msg
}

type StateLabelLanMessage struct {
	Label Label
}

func (o *StateLabelLanMessage) UnmarshalBinary(data []byte) error {
	o.Label = Label(bytes.TrimRight(data, "\x00"))

	return nil
}

func (o LanDeviceMessageBuilder) GetVersion() SendableLanMessage {
	return o.buildNormalMessageOfType(GetVersionType)
}

type StateVersionLanMessage struct {
	Vendor  uint32
	Product uint32
	Version uint32
}

func (o *StateVersionLanMessage) UnmarshalBinary(data []byte) error {
	// Vendor.
	o.Vendor = binary.LittleEndian.Uint32(data[:4])

	// Product.
	o.Product = binary.LittleEndian.Uint32(data[4:8])

	// Version.
	o.Version = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func (o LanDeviceMessageBuilder) GetInfo() SendableLanMessage {
	return o.buildNormalMessageOfType(GetInfoType)
}

type StateInfoLanMessage struct {
	Time     Time
	Uptime   uint64
	Downtime uint64
}

func (o *StateInfoLanMessage) UnmarshalBinary(data []byte) error {
	// Time.
	o.Time = Time(binary.LittleEndian.Uint64(data[:8]))

	// Uptime.
	o.Uptime = binary.LittleEndian.Uint64(data[8:16])

	// Downtime.
	o.Downtime = binary.LittleEndian.Uint64(data[16:])

	return nil
}

type AcknowledgementLanMessage struct{}

func (o *AcknowledgementLanMessage) UnmarshalBinary(data []byte) error {
	return nil
}

func (o LanDeviceMessageBuilder) GetLocation() SendableLanMessage {
	return o.buildNormalMessageOfType(GetLocationType)
}

type StateLocationLanMessage struct {
	Location  [16]byte
	Label     Label
	UpdatedAt Time
}

func (o *StateLocationLanMessage) UnmarshalBinary(data []byte) error {
	// Location.
	copy(o.Location[:], data[:16])

	// Label.
	o.Label = Label(bytes.TrimRight(data[16:48], "\x00"))

	// Updated at.
	o.UpdatedAt = Time(binary.LittleEndian.Uint64(data[48:]))

	return nil
}

func (o LanDeviceMessageBuilder) GetGroup() SendableLanMessage {
	return o.buildNormalMessageOfType(GetGroupType)
}

type StateGroupLanMessage struct {
	Group     [16]byte
	Label     Label
	UpdatedAt Time
}

func (o *StateGroupLanMessage) UnmarshalBinary(data []byte) error {
	// Group.
	copy(o.Group[:], data[:16])

	// Label.
	o.Label = Label(bytes.TrimRight(data[16:48], "\x00"))

	// Updated at.
	o.UpdatedAt = Time(binary.LittleEndian.Uint64(data[48:]))

	return nil
}

type EchoRequestLanMessage struct {
	Payload [64]byte
}

func (o EchoRequestLanMessage) MarshalBinary() ([]byte, error) {
	return o.Payload[:], nil
}

func (o LanDeviceMessageBuilder) EchoRequest(payload EchoRequestLanMessage) SendableLanMessage {
	msg := o.buildNormalMessageOfType(EchoRequestType)

	msg.Payload = payload
	msg.updateSize()

	return msg
}

type EchoResponseLanMessage struct {
	Payload [64]byte
}

func (o *EchoResponseLanMessage) UnmarshalBinary(data []byte) error {
	copy(o.Payload[:], data[:64])

	return nil
}

type HSBK struct {
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
}

func (o HSBK) MarshalBinary() (data []byte, _ error) {
	data = make([]byte, 8)

	// Hue.
	binary.LittleEndian.PutUint16(data[:2], o.Hue)

	// Saturation.
	binary.LittleEndian.PutUint16(data[2:4], o.Saturation)

	// Brightness.
	binary.LittleEndian.PutUint16(data[4:6], o.Brightness)

	if o.Kelvin < 2500 || o.Kelvin > 9000 {
		return nil, fmt.Errorf("color temperature %d out of range (2500..9000)", o.Kelvin)
	}

	// Kelvin.
	binary.LittleEndian.PutUint16(data[6:], o.Kelvin)

	return
}

func (o *HSBK) UnmarshalBinary(data []byte) error {
	// Hue.
	o.Hue = binary.LittleEndian.Uint16(data[:2])

	// Saturation.
	o.Saturation = binary.LittleEndian.Uint16(data[2:4])

	// Brightness.
	o.Brightness = binary.LittleEndian.Uint16(data[4:6])

	// Kelvin.
	o.Kelvin = binary.LittleEndian.Uint16(data[6:])

	return nil
}

func (o LanDeviceMessageBuilder) LightGet() SendableLanMessage {
	return o.buildNormalMessageOfType(LightGetType)
}

type LightSetColorLanMessage struct {
	Color    HSBK
	Duration uint32
}

func (o LightSetColorLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 13)

	// Color.
	color, err := o.Color.MarshalBinary()
	if err != nil {
		return
	}

	copy(data[1:9], color)

	// Duration.
	binary.LittleEndian.PutUint32(data[9:], o.Duration)

	return
}

func (o LanDeviceMessageBuilder) LightSetColor(payload LightSetColorLanMessage) SendableLanMessage {
	msg := o.buildNormalMessageOfType(LightSetColorType)

	msg.Payload = payload
	msg.updateSize()

	return msg
}

type LightStateLanMessage struct {
	Color HSBK
	Power PowerLevel
	Label Label
}

func (o *LightStateLanMessage) UnmarshalBinary(data []byte) error {
	// Color.
	err := o.Color.UnmarshalBinary(data[:8])
	if err != nil {
		return err
	}

	// Power.
	o.Power = PowerLevel(binary.LittleEndian.Uint16(data[8:10]))

	// Label.
	o.Label = Label(bytes.TrimRight(data[10:], "\x00"))

	return nil
}

func (o LanDeviceMessageBuilder) LightGetPower() SendableLanMessage {
	return o.buildNormalMessageOfType(LightGetPowerType)
}

type LightSetPowerLanMessage struct {
	Level    PowerLevel
	Duration uint32
}

func (o LightSetPowerLanMessage) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 6)

	// Level.
	level, err := o.Level.MarshalBinary()
	if err != nil {
		return
	}

	copy(data[:2], level)

	// Duration.
	binary.LittleEndian.PutUint32(data[2:], o.Duration)

	return
}

func (o LanDeviceMessageBuilder) LightSetPower(payload LightSetPowerLanMessage) SendableLanMessage {
	msg := o.buildNormalMessageOfType(LightSetPowerType)

	msg.Payload = payload
	msg.updateSize()

	return msg
}

type LightStatePowerLanMessage struct {
	Level PowerLevel
}

func (o *LightStatePowerLanMessage) UnmarshalBinary(data []byte) error {
	o.Level = PowerLevel(binary.LittleEndian.Uint16(data))

	return nil
}
