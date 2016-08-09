package controlifx

import (
	"math/rand"
	"net"
	"time"
	"errors"
	"reflect"
	"encoding"
)

const (
	NormalDiscoverTimeout = 250

	maxReadSize = LanHeaderSize+64
)

type (
	// Returns true when the given message is part of the expected
	// response(s). For example, the message may have to have a certain
	// type, and so false is returned for messages that don't match that
	// condition.
	Filter func(ReceivableLanMessage) bool

	// The first bool return should specify if the device should be
	// registered with the connector. The second bool return should be true
	// when discovery should continue, false otherwise. The message will be
	// the response to the discovery broadcast with a payload type of
	// StateServiceLanMessage, while the device will have the IP address and
	// MAC of the discovered device.
	DiscoverFilter func(ReceivableLanMessage, Device) (bool, bool)
)

type Device struct {
	Addr *net.UDPAddr
	Mac  uint64
}

type Connector struct {
	bcastAddr *net.UDPAddr
	conn      *net.UDPConn

	DiscoverTimeout int
	Devices         []Device
}

func (o *Connector) connect() error {
	if o.conn != nil {
		return nil
	}

	const PortStr = "56700"

	laddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4zero.String(), PortStr))
	if err != nil {
		return err
	}

	if o.bcastAddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4bcast.String(), PortStr)); err != nil {
		return err
	}

	o.conn, err = net.ListenUDP("udp", laddr)
	return err
}

func (o *Connector) send(addr *net.UDPAddr, msg SendableLanMessage) error {
	if err := o.connect(); err != nil {
		return err
	}

	if addr == nil {
		addr = o.bcastAddr
	}

	b, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = o.conn.WriteTo(b, addr)
	return err
}

func (o *Connector) bcastGetService() (uint32, error) {
	source := rand.Uint32()
	msg := LanDeviceMessageBuilder{source:source}.GetService()
	msg.Header.Frame.Tagged = true
	return source, o.send(nil, msg)
}

func (o *Connector) readMsg(filter Filter) (msg ReceivableLanMessage, raddr *net.UDPAddr, err error) {
	for {
		b := make([]byte, maxReadSize)
		var n int
		n, raddr, err = o.conn.ReadFromUDP(b)
		if err != nil {
			return
		}
		b = b[:n]

		msg = ReceivableLanMessage{}
		err = msg.UnmarshalBinary(b)
		if err == nil && filter(msg) {
			break
		}
	}
	return
}

func (o *Connector) DiscoverNDevices(n int) error {
	source, err := o.bcastGetService()
	if err != nil {
		return err
	}

	o.setReadDeadlineIfApplicable()

	for n > 0 {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			return err
		}
		o.Devices = append(o.Devices, Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		})
		n--
	}
	// Remove read deadline.
	return o.conn.SetDeadline(time.Time{})
}

func (o *Connector) DiscoverAllDevices() error {
	source, err := o.bcastGetService()
	if err != nil {
		return err
	}

	// Discovering all devices requires a timeout.
	if o.DiscoverTimeout == 0 {
		o.DiscoverTimeout = NormalDiscoverTimeout
	}
	o.setReadDeadlineIfApplicable()

	for {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				break
			}
			return err
		}
		o.Devices = append(o.Devices, Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		})
	}

	// Remove read deadline.
	return o.conn.SetDeadline(time.Time{})
}

func (o *Connector) DiscoverFilteredDevices(filter DiscoverFilter) error {
	source, err := o.bcastGetService()
	if err != nil {
		return err
	}

	o.setReadDeadlineIfApplicable()

	for {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			return err
		}
		d := Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		}
		register, cont := filter(msg, d)

		if register {
			o.Devices = append(o.Devices, d)
		}
		if !cont {
			break
		}
	}

	// Remove read deadline.
	return o.conn.SetDeadline(time.Time{})
}

func (o *Connector) setReadDeadlineIfApplicable() error {
	if o.DiscoverTimeout > 0 {
		return o.conn.SetDeadline(time.Now().Add(time.Duration(o.DiscoverTimeout)*time.Millisecond))
	}
	return nil
}

func (o *Connector) RemoveDevice(device Device) bool {
	for i, d := range o.Devices {
		if d.Mac == device.Mac {
			o.Devices = append(o.Devices[:i], o.Devices[i+1:]...)
			return true
		}
	}
	return false
}

func (o Connector) SendTo(device Device, msg SendableLanMessage) error {
	msg.Header.FrameAddress.Target = device.Mac
	return o.send(device.Addr, msg)
}

func (o Connector) SendToAll(msg SendableLanMessage) error {
	if len(o.Devices) == 0 {
		return errors.New("no devices; either none are connected or none were discovered")
	}

	for _, device := range o.Devices {
		if err := o.SendTo(device, msg); err != nil {
			return err
		}
	}
	return nil
}

func (o Connector) BlindSendToAll(msg SendableLanMessage) error {
	msg.Header.FrameAddress.Target = 0
	return o.send(nil, msg)
}

func (o Connector) GetResponseFrom(device Device, msg SendableLanMessage, filter Filter) (recMsg ReceivableLanMessage, err error) {
	source := rand.Uint32()
	msg.Header.Frame.Source = source
	msg.Header.Frame.Tagged = true
	msg.Header.FrameAddress.Target = device.Mac
	msg.Header.FrameAddress.ResRequired = true

	if err = o.send(device.Addr, msg); err != nil {
		return
	}

	recMsg, _, err = o.readMsg(func(msg ReceivableLanMessage) bool {
		return checkSourceAndFilter(msg, source, filter)
	})
	return
}

func (o Connector) GetResponseFromAll(msg SendableLanMessage, filter Filter) (recMsgs map[Device]ReceivableLanMessage, err error) {
	if len(o.Devices) == 0 {
		err = errors.New("no devices; either none are connected or none were discovered")
		return
	}

	source := rand.Uint32()
	msg.Header.Frame.Source = source
	msg.Header.Frame.Tagged = false
	msg.Header.FrameAddress.Target = 0
	msg.Header.FrameAddress.ResRequired = true

	err = o.send(nil, msg)
	if err != nil {
		return
	}

	recMsgs = make(map[Device]ReceivableLanMessage)
	n := len(o.Devices)
	for n > 0 {
		var recMsg ReceivableLanMessage
		recMsg, _, err = o.readMsg(func(msg ReceivableLanMessage) bool {
			return checkSourceAndFilter(msg, source, filter)
		})
		if err != nil {
			return
		}
		if device, err := o.findDevice(recMsg.Header.FrameAddress.Target); err == nil {
			recMsgs[device] = recMsg
			n--
		}
	}
	return
}

func (o Connector) findDevice(mac uint64) (Device, error) {
	for _, v := range o.Devices {
		if v.Mac == mac {
			return v, nil
		}
	}
	return Device{}, errors.New("device not found")
}

func TypeFilter(t encoding.BinaryUnmarshaler) Filter {
	return func(msg ReceivableLanMessage) bool {
		return reflect.TypeOf(msg.Payload).ConvertibleTo(reflect.TypeOf(t))
	}
}

func checkSourceAndFilter(msg ReceivableLanMessage, source uint32, filter Filter) bool {
	sourceOk := msg.Header.Frame.Source == source

	if filter != nil {
		return sourceOk && filter(msg)
	}
	return sourceOk
}
