package controlifx

import (
	"math/rand"
	"net"
	"time"
)

const (
	// NormalTimeout is a sane number of milliseconds to wait before timing out during discovery.
	NormalTimeout = 250

	maxReadSize = LanHeaderSize+64
)

type (
	// Filter returns false for messages that should not be further processed.
	Filter func(ReceivableLanMessage) bool

	// DiscoverFilter first returns whether or not the device should be registered and later returned after
	// discovery. The second return value specifies if discovery should continue if there's still time left.
	DiscoverFilter func(ReceivableLanMessage, Device) (register bool, cont bool)

	// Device is a LIFX device on the network.
	Device struct {
		// Addr is the remote address of the device.
		Addr *net.UDPAddr
		// Mac is the MAC address of the device.
		Mac  uint64
	}

	// Connection is the connection between the client and the network devices.
	Connection struct {
		bcastAddr *net.UDPAddr
		conn      *net.UDPConn
	}
)

func Connect() (o Connection, err error) {
	const PortStr = "56700"

	laddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4zero.String(), PortStr))
	if err != nil {
		return
	}

	if o.conn, err = net.ListenUDP("udp", laddr); err != nil {
		return
	}

	o.bcastAddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4bcast.String(), PortStr));

	return
}

func (o Connection) send(addr *net.UDPAddr, msg SendableLanMessage) error {
	b, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = o.conn.WriteTo(b, addr)

	return err
}

func (o Connection) readMsg(filter Filter) (msg ReceivableLanMessage, raddr *net.UDPAddr, err error) {
	for {
		b := make([]byte, maxReadSize)
		var n int
		n, raddr, err = o.conn.ReadFromUDP(b)
		if err != nil {
			return
		}
		b = b[:n]

		msg = ReceivableLanMessage{}
		if err = msg.UnmarshalBinary(b); err == nil && filter(msg) {
			break
		}
	}

	return
}

func (o Connection) Close() error {
	if o.conn != nil {
		return o.conn.Close()
	}

	return nil
}

// DiscoverDevices discovers as many devices as possible on the network within the timeout and filters devices.
func (o Connection) DiscoverDevices(timeout int, filter DiscoverFilter) (devices []Device, err error) {
	getServiceMsg := GetService()
	getServiceMsg.Header.Frame.Source = rand.Uint32()

	if err = o.send(o.bcastAddr, getServiceMsg); err != nil {
		return
	}

	o.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout)*time.Millisecond))

	for {
		recMsg, raddr, err := o.readMsg(func(recMsg ReceivableLanMessage) bool {
			payload, ok := recMsg.Payload.(*StateServiceLanMessage)

			return recMsg.Header.Frame.Source == getServiceMsg.Header.Frame.Source &&
				ok && payload.Service == 1
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				err = nil
			}

			break
		}

		d := Device{
			Addr: raddr,
			Mac:  recMsg.Header.FrameAddress.Target,
		}

		if filter == nil {
			devices = append(devices, d)
		} else {
			register, cont := filter(recMsg, d)
			if register {
				devices = append(devices, d)
			}
			if !cont {
				break
			}
		}
	}

	// Remove read deadline.
	o.conn.SetReadDeadline(time.Time{})

	return
}

// DiscoverAllDevices discovers as many devices as possible on the network within the timeout.
func (o Connection) DiscoverAllDevices(timeout int) ([]Device, error) {
	return o.DiscoverDevices(timeout, nil)
}

// SendTo sends the message to the devices without expecting responses.
func (o Connection) SendTo(msg SendableLanMessage, devices []Device) error {
	// Possible bug in LIFX protocol. Workaround here:
	switch msg.Header.ProtocolHeader.Type {
	case SetPowerType:
	case LightSetPowerType:
	default:
		msg.Header.Frame.Tagged = true
	}

	for _, d := range devices {
		msg.Header.FrameAddress.Target = d.Mac

		if err := o.send(d.Addr, msg); err != nil {
			return err
		}
	}

	return nil
}

// SendToAll sends the message to all devices on the network without expecting responses.
func (o Connection) SendToAll(msg SendableLanMessage) error {
	msg.Header.Frame.Tagged = false
	msg.Header.FrameAddress.Target = 0

	return o.send(o.bcastAddr, msg)
}

// SendToAndGet sends the message to the devices, filters the responses, and builds a mapping between a responding
// device and its response.
func (o Connection) SendToAndGet(msg SendableLanMessage, devices []Device, filter Filter) (recMsgs map[Device]ReceivableLanMessage, err error) {
	msg.Header.Frame.Source = rand.Uint32()
	msg.Header.FrameAddress.ResRequired = true

	if err = o.SendToAll(msg); err != nil {
		return
	}

	recMsgs = make(map[Device]ReceivableLanMessage)

	for len(devices) > 0 {
		recMsg, _, err := o.readMsg(func(recMsg ReceivableLanMessage) bool {
			return checkSourceAndFilter(recMsg, msg.Header.Frame.Source, filter)
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				err = nil
			}

			break
		}

		for i, d := range devices {
			if d.Mac == recMsg.Header.FrameAddress.Target {
				recMsgs[d] = recMsg
				devices = append(devices[:i], devices[i+1:]...)
				break
			}
		}
	}

	return
}

// SendToAllAndGet sends the message to all devices on the network, filters the responses, and builds a mapping between
// a responding device and its response.
func (o Connection) SendToAllAndGet(timeout int, msg SendableLanMessage, filter Filter) (recMsgs map[Device]ReceivableLanMessage, err error) {
	msg.Header.Frame.Source = rand.Uint32()
	msg.Header.FrameAddress.ResRequired = true

	if err = o.SendToAll(msg); err != nil {
		return
	}

	o.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout)*time.Millisecond))

	recMsgs = make(map[Device]ReceivableLanMessage)

	for {
		recMsg, raddr, err := o.readMsg(func(recMsg ReceivableLanMessage) bool {
			return checkSourceAndFilter(recMsg, msg.Header.Frame.Source, filter)
		})
		if err != nil {
			if err.(net.Error).Timeout() {
				err = nil
			}

			break
		}

		d := Device{
			Addr:raddr,
			Mac:recMsg.Header.FrameAddress.Target,
		}

		recMsgs[d] = recMsg
	}

	return
}

// TypeFilter filters out responses that do not have the payload type.
func TypeFilter(t uint16) Filter {
	return func(msg ReceivableLanMessage) bool {
		return msg.Header.ProtocolHeader.Type == t
	}
}

func checkSourceAndFilter(msg ReceivableLanMessage, source uint32, filter Filter) bool {
	sourceOk := msg.Header.Frame.Source == source

	if filter != nil {
		return sourceOk && filter(msg)
	}

	return sourceOk
}
