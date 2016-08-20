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
	// NormalDiscoveryTimeout A sane discover timeout specified in milliseconds.
	NormalDiscoverTimeout = 250

	maxReadSize = LanHeaderSize+64
)

type (
	// Filter is responsible for filtering out extraneous network responses
	// that should not be further processed. It returns true when the given
	// message should be further processed. For example, the message may
	// have to have a certain type, and so false should be returned for
	// messages that don't match that condition.
	Filter func(ReceivableLanMessage) bool

	// DiscoverFilter is responsible for controlling which devices are to be
	// registered and whether more devices should be discovered. The first
	// return value should specify if the device should be registered with
	// the connector. The second return value should be true when discovery
	// should continue, or false otherwise. The message will be the response
	// to the discovery broadcast with a payload type of
	// StateServiceLanMessage, while the device will be the not yet
	// registered responding device.
	DiscoverFilter func(ReceivableLanMessage, Device) (bool, bool)
)

// Device is a LIFX device on the LAN.
type Device struct {
	// Addr is the remote address of the device.
	Addr *net.UDPAddr
	// Mac is the MAC address of the device.
	Mac  uint64
}

// Connector is the connection between the client and network devices.
type Connector struct {
	bcastAddr *net.UDPAddr
	conn      *net.UDPConn

	// DiscoverTimeout is the maximum number of milliseconds to wait to
	// discover devices. A zero value represents no timeout, however a sane
	// one will be used regardless if discovering via DiscoverAllDevices.
	DiscoverTimeout int
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

func (o Connector) readMsg(filter Filter) (msg ReceivableLanMessage, raddr *net.UDPAddr, err error) {
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
			return
		}
	}
	return
}

// DiscoverNDevices discovers n devices on the network, returning as soon as n
// devices respond or when DiscoverTimeout is reached, whichever comes first.
// If DiscoverTimeout is the zero value, this function will not time out.
func (o *Connector) DiscoverNDevices(n int) ([]Device, error) {
	devices := make([]Device, n)
	source, err := o.bcastGetService()
	if err != nil {
		return devices, err
	}
	o.conn.SetReadDeadline(o.getReadDeadlineIfApplicable())
	var i int
	for i < n {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				break
			}
			return devices, err
		}
		devices[i] = Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		}
		i++
	}
	// Remove read deadline.
	return devices, o.conn.SetDeadline(time.Time{})
}

// DiscoverAllDevices discovers as many devices as possible until
// DiscoverTimeout is reached. If DiscoverTimeout is the zero value, a sane
// default will be used (250 ms).
func (o *Connector) DiscoverAllDevices() ([]Device, error) {
	var devices []Device
	source, err := o.bcastGetService()
	if err != nil {
		return devices, err
	}

	// Discovering all devices requires a timeout.
	readDeadline := o.getReadDeadlineIfApplicable()
	if readDeadline.IsZero() {
		readDeadline = time.Now().Add(NormalDiscoverTimeout*time.Millisecond)
	}
	o.conn.SetReadDeadline(readDeadline)

	for {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				break
			}
			return devices, err
		}
		devices = append(devices, Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		})
	}

	// Remove read deadline.
	return devices, o.conn.SetDeadline(time.Time{})
}

// DiscoverFilteredDevices discovers devices until DiscoverTimeout is reached or
// the filter's second return value is false, whichever comes first. If
// DiscoverTimeout is the zero value, this function will not time out.
func (o *Connector) DiscoverFilteredDevices(filter DiscoverFilter) ([]Device, error) {
	var devices []Device
	source, err := o.bcastGetService()
	if err != nil {
		return devices, err
	}
	o.conn.SetReadDeadline(o.getReadDeadlineIfApplicable())
	for {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				break
			}
			return devices, err
		}
		d := Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		}
		register, cont := filter(msg, d)
		if register {
			devices = append(devices, d)
		}
		if !cont {
			break
		}
	}

	// Remove read deadline.
	return devices, o.conn.SetDeadline(time.Time{})
}

func (o Connector) getReadDeadlineIfApplicable() time.Time {
	if o.DiscoverTimeout > 0 {
		return time.Now().Add(time.Duration(o.DiscoverTimeout)*time.Millisecond)
	}
	return time.Time{}
}

// SendTo sends the message to each device, not expecting responses.
func (o *Connector) SendTo(msg SendableLanMessage, devices []Device) error {
	msg.Header.Frame.Tagged = true
	ch := make(chan error)
	for _, v := range devices {
		go func(msg SendableLanMessage, device Device) {
			msg.Header.FrameAddress.Target = device.Mac
			ch <- o.send(device.Addr, msg)
		}(msg, v)
	}
	var i int
	for i < len(devices) {
		select {
		case err := <-ch:
			i++
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SendToAll sends the message to all devices on the network, not expecting
// responses.
func (o *Connector) SendToAll(msg SendableLanMessage) error {
	return o.send(nil, msg)
}

// SendToAndGet sends the message to the devices, filtering each response and
// returning a mapping between the responding device and its response.
func (o *Connector) SendToAndGet(msg SendableLanMessage, filter Filter, devices []Device) (recMsgs map[Device]ReceivableLanMessage, err error) {
	source := rand.Uint32()
	msg.Header.Frame.Source = source
	msg.Header.FrameAddress.ResRequired = true
	if err = o.SendTo(msg, devices); err != nil {
		return
	}
	recMsgs = make(map[Device]ReceivableLanMessage)
	n := len(devices)
	for n > 0 {
		var recMsg ReceivableLanMessage
		recMsg, _, err = o.readMsg(func(msg ReceivableLanMessage) bool {
			return checkSourceAndFilter(msg, source, filter)
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				err = nil
			}
			break
		}
		if device, err := o.findDevice(recMsg.Header.FrameAddress.Target, devices); err == nil {
			recMsgs[device] = recMsg
			n--
		}
	}
	return
}

// SendToAndGet sends the message to all devices on the network, filtering each
// response and returning a mapping between the responding device and its
// response.
func (o *Connector) SendToAllAndGet(msg SendableLanMessage, filter Filter) (recMsgs map[Device]ReceivableLanMessage, err error) {
	source := rand.Uint32()
	msg.Header.Frame.Source = source
	msg.Header.Frame.Tagged = false
	if err = o.SendToAll(msg); err != nil {
		return
	}
	recMsgs = make(map[Device]ReceivableLanMessage)

	// Receiving responses from all devices requires a timeout.
	readDeadline := o.getReadDeadlineIfApplicable()
	if readDeadline.IsZero() {
		readDeadline = time.Now().Add(NormalDiscoverTimeout*time.Millisecond)
	}
	o.conn.SetReadDeadline(readDeadline)

	for {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			return checkSourceAndFilter(msg, source, filter)
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				return recMsgs, nil
			}
			return recMsgs, err
		}
		device := Device{
			Addr: raddr,
			Mac: msg.Header.FrameAddress.Target,
		}
		recMsgs[device] = msg
	}
	return
}

func (o Connector) findDevice(mac uint64, devices []Device) (Device, error) {
	for _, v := range devices {
		if v.Mac == mac {
			return v, nil
		}
	}
	return Device{}, errors.New("device not found")
}

// TypeFilter filters out responses that do not have a payload of the given type
// t. Since assuring that the payload of a response is a certain type is such a
// common task when receiving responses, this func has been provided for
// convenience.
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
