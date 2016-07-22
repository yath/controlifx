package controlifx

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"
	"errors"
)

const maxReadSize = LanHeaderSize+64

type Filter func(ReceivableLanMessage) bool

type device struct {
	addr *net.UDPAddr
	mac  uint64
}

type Connector struct {
	bcastAddr *net.UDPAddr
	conn      *net.UDPConn

	Devices []device
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
	source, err := randomSource()
	if err != nil {
		return source, err
	}

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

	for n > 0 {
		msg, raddr, err := o.readMsg(func(msg ReceivableLanMessage) bool {
			payload, ok := msg.Payload.(*StateServiceLanMessage)
			return msg.Header.Frame.Source == source && ok && payload.Service == 1
		})
		if err != nil {
			return err
		}
		o.Devices = append(o.Devices, device{
			addr: raddr,
			mac: msg.Header.FrameAddress.Target,
		})
		n--
	}
	return nil
}

func (o *Connector) DiscoverAllDevices(timeout int) error {
	source, err := o.bcastGetService()
	if err != nil {
		return err
	}

	if err := o.conn.SetDeadline(time.Now().Add(time.Duration(timeout)*time.Millisecond)); err != nil {
		return err
	}

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
		o.Devices = append(o.Devices, device{
			addr: raddr,
			mac: msg.Header.FrameAddress.Target,
		})
	}

	// Remove read deadline.
	return o.conn.SetDeadline(time.Time{})
}

func (o Connector) SendTo(device device, msg SendableLanMessage) error {
	msg.Header.FrameAddress.Target = device.mac
	return o.send(device.addr, msg)
}

func (o Connector) SendToAll(msg SendableLanMessage) error {
	msg.Header.FrameAddress.Target = 0
	return o.send(nil, msg)
}

func (o Connector) GetResponseFrom(device device, msg SendableLanMessage, filter Filter) (recMsg ReceivableLanMessage, err error) {
	source, err := randomSource()
	if err != nil {
		return
	}

	msg.Header.Frame.Source = source
	msg.Header.Frame.Tagged = true
	msg.Header.FrameAddress.Target = device.mac
	msg.Header.FrameAddress.ResRequired = true

	if err = o.send(device.addr, msg); err != nil {
		return
	}

	recMsg, _, err = o.readMsg(func(msg ReceivableLanMessage) bool {
		return checkSourceAndFilter(msg, source, filter)
	})
	return
}

func (o Connector) GetResponseFromAll(msg SendableLanMessage, filter Filter) (recMsgs map[device]ReceivableLanMessage, err error) {
	if len(o.Devices) == 0 {
		err = errors.New("no devices; either none are connected or none were discovered")
		return
	}

	source, err := randomSource()
	if err != nil {
		return
	}

	msg.Header.Frame.Source = source
	msg.Header.Frame.Tagged = true
	msg.Header.FrameAddress.Target = 0
	msg.Header.FrameAddress.ResRequired = true

	err = o.send(nil, msg)
	if err != nil {
		return
	}

	recMsgs = make(map[device]ReceivableLanMessage)
	n := len(o.Devices)
	for n > 0 {
		var recMsg ReceivableLanMessage
		recMsg, _, err = o.readMsg(func(msg ReceivableLanMessage) bool {
			return checkSourceAndFilter(msg, source, filter)
		})
		if err != nil {
			return
		}
		var device device
		device, err = o.findDevice(recMsg.Header.FrameAddress.Target)
		if err != nil {
			return
		}
		recMsgs[device] = recMsg
		n--
	}
	return
}

func (o Connector) findDevice(mac uint64) (device, error) {
	for _, v := range o.Devices {
		if v.mac == mac {
			return v, nil
		}
	}
	return device{}, errors.New("device not found")
}

func checkSourceAndFilter(msg ReceivableLanMessage, source uint32, filter Filter) bool {
	sourceOk := msg.Header.Frame.Source == source

	if filter != nil {
		return sourceOk && filter(msg)
	}
	return sourceOk
}

func randomSource() (uint32, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	return binary.BigEndian.Uint32(b), err
}
