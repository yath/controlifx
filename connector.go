package controlifx

import (
	"net"
	_time "time"
	"math/rand"
	"encoding/binary"
	"github.com/siddontang/go/log"
)

const MaxReadSize = LanHeaderSize + 64

type Device struct {
	addr *net.UDPAddr
	mac uint64
}

type Connector struct {
	ins []chan receivablePacket
	out chan<- sendablePacket

	Devices []Device
}

type receivablePacket struct {
	ReceivableLanMessage

	Device Device
}

type sendablePacket struct {
	SendableLanMessage

	device *Device
}

func (o *Connector) getIn() (int, <-chan receivablePacket) {
	for i, v := range o.ins {
		if v == nil {
			o.ins[i] = make(chan receivablePacket)
			return i, o.ins[i]
		}
	}

	i := len(o.ins)
	o.ins = append(o.ins, make(chan receivablePacket))
	return i, o.ins[i]
}

func (o *Connector) doneWithIn(i int) {
	if i == len(o.ins)-1 {
		o.ins = o.ins[:i]
	} else {
		o.ins[i] = nil
	}
}

func (o *Connector) Connect() error {
	const PortStr = "56700"

	laddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4zero.String(), PortStr))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", laddr)

	out := make(chan sendablePacket)
	o.out = out

	bcastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(net.IPv4bcast.String(), PortStr))
	if err != nil {
		return err
	}

	// In.
	go func() {
		for {
			b := make([]byte, MaxReadSize)
			n, addr, err := conn.ReadFromUDP(b)
			if err != nil {
				log.Error(err)
				continue
			}
			b = b[:n]

			msg := ReceivableLanMessage{}
			err = msg.UnmarshalBinary(b)
			if err != nil {
				continue
			}

			for _, v := range o.ins {
				if v != nil {
					v <- receivablePacket{
						msg,
						Device{
							addr:addr,
							mac:msg.header.frameAddress.Target,
						},
					}
				}
			}
		}
	}()

	// Out.
	go func() {
		for {
			packet := <-out

			b, err := packet.MarshalBinary()
			if err != nil {
				log.Error(err)
				continue
			}

			var dest *net.UDPAddr

			if packet.device == nil {
				dest = bcastAddr
			} else {
				dest = packet.device.addr
				packet.header.frameAddress.Target = packet.device.mac
			}

			if err != nil {
				log.Error(err)
				continue
			}

			_, err = conn.WriteTo(b, dest)
			if err != nil {
				log.Error(err)
			}
		}

		close(out)
	}()

	return err
}

func (o *Connector) FindAllDevices() error {
	const DiscoveryTimeout = 2

	source := RandomSource()

	msg := LanDeviceMessageBuilder{
		Source:source,
	}.GetService()
	msg.header.frame.Tagged = true

	o.out <- sendablePacket{msg, nil}
	inIndex, in := o.getIn()

	timeout := _time.NewTimer(DiscoveryTimeout * 1e9)

	RECEIVE_LOOP:
	for {
		select {
		case <-timeout.C:
			break RECEIVE_LOOP
		case msg := <-in:
			if serviceMsg, ok := msg.payload.(*StateServiceLanMessage); ok && msg.header.frame.Source == source && serviceMsg.Service == 1 {
				o.Devices = append(o.Devices, msg.Device)
			}
		}
	}

	log.Debugf("Registered %d devices\n", len(o.Devices))

	o.doneWithIn(inIndex)
	return nil
}

type Filter func(ReceivableLanMessage) bool

func (o *Connector) SendMessageToAll(msg SendableLanMessage) {
	o.out <- sendablePacket{msg, nil}
}

func (o *Connector) SendMessageTo(device *Device, msg SendableLanMessage) {
	o.out <- sendablePacket{msg, device}
}

func (o *Connector) GetResponseFromAll(msg SendableLanMessage, filter Filter) <-chan receivablePacket {
	inIndex, in := o.getIn()

	o.SendMessageToAll(msg)

	ch := make(chan receivablePacket)

	go func() {
		var good int

		for good < len(o.Devices) {
			msg := <-in
			if filter(msg.ReceivableLanMessage) {
				ch <- msg
				good++
			}
		}

		close(ch)
		o.doneWithIn(inIndex)
	}()

	return ch
}

func (o *Connector) GetResponseFrom(device *Device, msg SendableLanMessage, filter Filter) <-chan receivablePacket {
	inIndex, in := o.getIn()

	o.SendMessageTo(device, msg)

	ch := make(chan receivablePacket)

	go func() {
		for {
			msg := <-in
			if filter(msg.ReceivableLanMessage) {
				ch <- msg
				break
			}
		}

		close(ch)
		o.doneWithIn(inIndex)
	}()

	return ch
}

func RandomSource() uint32 {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return binary.BigEndian.Uint32(b)
}
