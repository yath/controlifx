package main

import (
	"github.com/bionicrm/controlifx"
	"log"
)

func main() {
	conn := controlifx.Connector{}

	if err := conn.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := conn.FindAllDevices(1); err != nil {
		log.Fatal(err)
	}

	builder := controlifx.LanDeviceMessageBuilder{}
	payload := controlifx.LightSetColorLanMessage{
		Color:controlifx.HSBK{
			Hue:0xffff,
			Saturation:0xffff,
			Brightness:0xffff,
			Kelvin:3500,
		},
		Duration:10*1000,
	}
	msg := builder.LightSetColor(payload)

	conn.SendMessageToAll(msg)
}
