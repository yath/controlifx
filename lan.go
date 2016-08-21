package controlifx

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"time"
	"math"
)

const (
	// MessageRate is the maximum recommended number of messages a device should
	// receive in a second.
	MessageRate = 20

	// The LanHeaderSize is the size of the header for a LAN message in bytes.
	LanHeaderSize = 36
)

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
	payload, ok := receivableTypeToPayload[o.Header.ProtocolHeader.Type]
	if !ok {
		return fmt.Errorf("cannot create new payload of type %d; is it binary decodable?",
			o.Header.ProtocolHeader.Type)
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

	// Little endian.
	putUint48 := func (b []byte, v uint64) {
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
	}

	// Target.
	if o.Target > 0xffffffffffff {
		binary.LittleEndian.PutUint64(data[:8], o.Target)
	} else {
		putUint48(data[:6], o.Target)
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
	// Little endian.
	uint48 := func(b []byte) uint64 {
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40
	}

	// Target.
	if data[7] | data[8] == 0 {
		o.Target = uint48(data[:6])
	} else {
		o.Target = binary.LittleEndian.Uint64(data[:8])
	}

	// 0000 00??

	// AckRequired.
	o.AckRequired = ((data[14]>>1)&0x01 == 1)

	// ResRequired.
	o.ResRequired = (data[14]&0x01 == 1)

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

	data = append([]byte(o), make([]byte, Size-len(o))...)

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

const (
	// Sendable types.
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

const (
	// Receivable types.
	StateServiceType      = 3
	StateHostInfoType     = 13
	StateHostFirmwareType = 15
	StateWifiInfoType     = 17
	StateWifiFirmwareType = 19
	StatePowerType        = 22
	StateLabelType        = 25
	StateVersionType      = 33
	StateInfoType         = 35
	AcknowledgementType   = 45
	StateLocationType     = 50
	StateGroupType        = 53
	EchoResponseType      = 59
	LightStateType        = 107
	LightStatePowerType   = 118
)

var receivableTypeToPayload = map[uint16]encoding.BinaryUnmarshaler{
	StateServiceType:      &StateServiceLanMessage{},
	StateHostInfoType:     &StateHostInfoLanMessage{},
	StateHostFirmwareType: &StateHostFirmwareLanMessage{},
	StateWifiInfoType:     &StateWifiInfoLanMessage{},
	StateWifiFirmwareType: &StateWifiFirmwareLanMessage{},
	StatePowerType:        &StatePowerLanMessage{},
	StateLabelType:        &StateLabelLanMessage{},
	StateVersionType:      &StateVersionLanMessage{},
	StateInfoType:         &StateInfoLanMessage{},
	AcknowledgementType:   &AcknowledgementLanMessage{},
	StateLocationType:     &StateLocationLanMessage{},
	StateGroupType:        &StateGroupLanMessage{},
	EchoResponseType:      &EchoResponseLanMessage{},
	LightStateType:        &LightStatePowerLanMessage{},
	LightStatePowerType:   &LightStatePowerLanMessage{},
}

func createSendableLanMessage(t uint16) SendableLanMessage {
	return SendableLanMessage{
		Header:LanHeader{
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:t,
			},
		},
	}
}

func GetService() SendableLanMessage {
	return createSendableLanMessage(GetServiceType)
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

func GetHostInfo() SendableLanMessage {
	return createSendableLanMessage(GetHostInfoType)
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

func GetHostFirmware() SendableLanMessage {
	return createSendableLanMessage(GetHostFirmwareType)
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

func GetWifiInfo() SendableLanMessage {
	return createSendableLanMessage(GetWifiInfoType)
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

func GetWifiFirmware() SendableLanMessage {
	return createSendableLanMessage(GetWifiFirmwareType)
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

func GetPower() SendableLanMessage {
	return createSendableLanMessage(GetPowerType)
}

type SetPowerLanMessage struct {
	Level PowerLevel
}

func (o SetPowerLanMessage) MarshalBinary() ([]byte, error) {
	return o.Level.MarshalBinary()
}

func SetPower(payload SetPowerLanMessage) SendableLanMessage {
	msg := createSendableLanMessage(SetPowerType)
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

func GetLabel() SendableLanMessage {
	return createSendableLanMessage(GetLabelType)
}

type SetLabelLanMessage struct {
	Label Label
}

func (o SetLabelLanMessage) MarshalBinary() ([]byte, error) {
	return o.Label.MarshalBinary()
}

func SetLabel(payload SetLabelLanMessage) SendableLanMessage {
	msg := createSendableLanMessage(SetLabelType)
	msg.Payload = payload

	msg.updateSize()

	return msg

	return msg
}

type StateLabelLanMessage struct {
	Label Label
}

func (o *StateLabelLanMessage) UnmarshalBinary(data []byte) error {
	o.Label = Label(bytes.TrimRight(data, "\x00"))

	return nil
}

func GetVersion() SendableLanMessage {
	return createSendableLanMessage(GetVersionType)
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

func GetInfo() SendableLanMessage {
	return createSendableLanMessage(GetInfoType)
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

func GetLocation() SendableLanMessage {
	return createSendableLanMessage(GetLocationType)
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

func GetGroup() SendableLanMessage {
	return createSendableLanMessage(GetGroupType)
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

func EchoRequest(payload EchoRequestLanMessage) SendableLanMessage {
	msg := createSendableLanMessage(EchoRequestType)
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

func LightGet() SendableLanMessage {
	return createSendableLanMessage(LightGetType)
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

func LightSetColor(payload LightSetColorLanMessage) SendableLanMessage {
	msg := createSendableLanMessage(LightSetColorType)
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

func LightGetPower() SendableLanMessage {
	return createSendableLanMessage(LightGetPowerType)
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

func LightSetPower(payload LightSetPowerLanMessage) SendableLanMessage {
	msg := createSendableLanMessage(LightSetPowerType)
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
