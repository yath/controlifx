# Controlifx
Client side API for LIFX device control.

**Projects built with Controlifx:**
- [clifx](https://github.com/golifx/clifx) &ndash; command-line interface for LIFX device control
- [emulifx](https://github.com/golifx/emulifx) &ndash; LIFX device emulator
- [implifx](https://github.com/golifx/implifx) &ndash; server side API for LIFX device implementations

**Resources:**
- [LIFX LAN protocol](https://lan.developer.lifx.com/)

**Contents:**
- [Installation](#installation)
- [Getting Started](#getting-started)

### Installation
Just run `go get -u gopkg.in/bionicrm/controlifx.v1` to get the latest version.

### Getting Started
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

Finally, the `TypeFilter(...)` will assure that we only process messages that have the type we're expecting, since other clients on the network may cause devices to emit other types of responses. This assures we don't try to process a message telling us information about its version when we just wanted to know its label.
```go
recMsgs, err := conn.SendToAllAndGet(controlifx.NormalTimeout, msg,
		controlifx.TypeFilter(controlifx.StateLabelType))
if err != nil {
	log.Fatalln(err)
}
```
We now have all of our responses stored in `recMsgs`, which is a mapping between a responding device and its response. The type cast to `*controlifx.StateLabelLanMessage` will always succeed because we know for sure, thanks to our `TypeFilter(...)` above, that the message is in fact a StateLabel response.
```go
for device, recMsg := range recMsgs {
	payload := recMsg.Payload.(*controlifx.StateLabelLanMessage)

	log.Printf("Received StateLabel response from %s: '%s'\n",
		device.Addr.String(), payload.Label)
}
```
**Here is the completed example:**
```go
package main

import (
	"gopkg.in/bionicrm/controlifx.v1"
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
