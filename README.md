# Controlifx
Client side API for LIFX device control.

**Projects built with Controlifx:**
- [Clifx](https://github.com/lifx-tools/clifx) &ndash; command-line interface for LIFX device control
- [Emulifx](https://github.com/lifx-tools/emulifx) &ndash; LIFX device emulator
- [Implifx](https://github.com/lifx-tools/implifx) &ndash; server side API for LIFX device implementations

**Resources:**
- [LIFX LAN protocol](https://lan.developer.lifx.com/)

**Contents:**
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Examples](#examples)
  - [Changing colors](#changing-colors) &ndash; change the color of the bulbs
- [Additional Help](#additional-help)

## Installation
Just run `go get -u gopkg.in/lifx-tools/controlifx.v1` to get the latest version.

## Getting Started
You'll always start off by opening up a UDP socket for sending and receiving messages:

```go
conn, err := controlifx.Connect()
if err != nil {
    log.Fatalln(err)
}
defer conn.Close()
```

Next you'll need to create a message to send to LIFX devices on the LAN. All of the messages can be found and described on the official [device messages](https://lan.developer.lifx.com/docs/device-messages) and [light messages](https://lan.developer.lifx.com/docs/light-messages) docs. Remember that some messages are only sent by devices to the client (usually those that start with "State"), and so they cannot be sent from the client, only received. Here, `GetLabel()` will return a new device message that will eventually make the LIFX devices emit a `StateLabel` response.

```go
msg := controlifx.GetLabel()
```
Now we need to send the message and wait for responses. `SendToAllAndGet(...)` will emit a message on the network that is received by all connected LIFX devices. The timeout value is how long to wait for the devices to respond, since we don't know how many of them will actually respond to the message (we'll talk about alternative methods later). 

Finally, the `TypeFilter(...)` will assure that we only process messages that have the type we're expecting, since other clients on the network may cause devices to emit other types of responses. This ensures we don't try to process a message telling us information about its version when we just wanted to know its label.

```go
recMsgs, err := conn.SendToAllAndGet(controlifx.NormalTimeout, msg,
	controlifx.TypeFilter(controlifx.StateLabelType))
if err != nil {
	log.Fatalln(err)
}
```

We now have all of our responses stored in `recMsgs`, which is a mapping between a responding device and its response. The type assertion of `*controlifx.StateLabelLanMessage` will always succeed because we know for sure, thanks to our `TypeFilter(...)` above, that the message is in fact a StateLabel response.

```go
for device, recMsg := range recMsgs {
	payload := recMsg.Payload.(*controlifx.StateLabelLanMessage)

	log.Printf("Received StateLabel response from %s: '%s'\n",
		device.Addr.String(), payload.Label)
}
```

**Completed example:**
```go
package main

import (
	"gopkg.in/lifx-tools/controlifx.v1"
	"log"
)

func main() {
	conn, err := controlifx.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	msg := controlifx.GetLabel()

	recMsgs, err := conn.SendToAllAndGet(controlifx.NormalTimeout, msg,
		controlifx.TypeFilter(controlifx.StateLabelType))
	if err != nil {
		log.Fatalln(err)
	}

	for device, recMsg := range recMsgs {
		payload := recMsg.Payload.(*controlifx.StateLabelLanMessage)

		log.Printf("Received StateLabel response from %s: '%s'\n",
			device.Addr.String(), payload.Label)
	}
}

```

**Example output:**
```
2016/08/25 13:13:48 Received StateLabel response from 10.0.0.23:56700: 'Floor'
2016/08/25 13:13:48 Received StateLabel response from 10.0.0.111:56700: 'Nightstand'
2016/08/25 13:13:48 Received StateLabel response from 10.0.0.132:56700: 'Closet'
```

## Examples
#### Changing colors
You'll undoubtedly want to change the light color of your LIFX bulbs at some point. In this example, we have to give the devices a payload so that they know what color we want them set to.

Afterwards, we send the message to all devices on the LAN. However, unlike in the [Getting Started](#getting-started) section, we use `SendToAll(...)` to emit the message on the network, ignoring responses. The method returns as soon as the message is written out to the network, unlike `SendToAllAndGet(...)`, which has to block until it receives responses. 

```go
package main

import (
	"gopkg.in/lifx-tools/controlifx.v1"
	"log"
)

func main() {
	// Open the UDP socket.
	conn, err := controlifx.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Create the payload to send to the LIFX devices.
	payload := controlifx.LightSetColorLanMessage{
		Color: controlifx.HSBK{
			Hue:        0xffff / 2,
			Saturation: 0xffff,
			Brightness: 0xffff,
			Kelvin:     3500,
		},
		Duration: 1500,
	}

	// Create the message with the payload that we're going to emit.
	msg := controlifx.LightSetColor(payload)

	// Send the message to all devices on the LAN, ignoring any responses.
	if err := conn.SendToAll(msg); err != nil {
		log.Fatalln(err)
	}
}

```

## Additional Help
Visit [#lifx-tools](http://webchat.freenode.net?randomnick=1&channels=%23lifx-tools&prompt=1) on chat.freenode.net to get help, ask questions, or discuss ideas.
