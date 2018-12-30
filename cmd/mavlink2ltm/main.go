package main

import (
	"flag"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mavlink"
	common "gobot.io/x/gobot/platforms/mavlink/common"
	"log"
)

func main() {
	var mavlinkUDP, serialPort string
	var serialBaud uint

	flag.StringVar(&mavlinkUDP, "mavlink", "127.0.0.1:14550", "Example: -mavlink 127.0.0.1:14550")
	flag.StringVar(&serialPort, "serial", "", "Example: -serial /dev/cu.SLAB_USBtoUART")
	flag.UintVar(&serialBaud, "baud", 2400, "Example: -baud 2400")
	flag.Parse()

	if serialPort == "" {
		panic("Set -serial (-help for help)")
	}

	fmt.Println(fmt.Sprintf("Converting from %s to %s at %dbaud", mavlinkUDP, serialPort, serialBaud))

	adaptor := mavlink.NewUDPAdaptor(mavlinkUDP)
	iris := mavlink.NewDriver(adaptor)

	// Set up options.
	options := serial.OpenOptions{
		PortName:        serialPort,
		BaudRate:        serialBaud,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	work := func() {
		iris.On(mavlink.MessageEvent, func(data interface{}) {
			if data.(common.MAVLinkMessage).Id() == 24 {
				message := data.(*common.GpsRawInt)
				fmt.Println(fmt.Sprintf("Sent %d, %d ALT %d", message.LAT, message.LON, message.ALT))

				crc := byte(0x00)
				var payload = make([]byte, 14)

				// LAT
				payload[0] = byte((message.LAT >> (8 * 0)) & 0xff)
				payload[1] = byte((message.LAT >> (8 * 1)) & 0xff)
				payload[2] = byte((message.LAT >> (8 * 2)) & 0xff)
				payload[3] = byte((message.LAT >> (8 * 3)) & 0xff)

				// LON
				payload[4] = byte((message.LON >> (8 * 0)) & 0xff)
				payload[5] = byte((message.LON >> (8 * 1)) & 0xff)
				payload[6] = byte((message.LON >> (8 * 2)) & 0xff)
				payload[7] = byte((message.LON >> (8 * 3)) & 0xff)

				// TODO: ground speed
				payload[8] = 0

				// ALT
				payload[9] = byte((message.ALT >> (8 * 0)) & 0xff)
				payload[10] = byte((message.ALT >> (8 * 1)) & 0xff)
				payload[11] = byte((message.ALT >> (8 * 2)) & 0xff)
				payload[12] = byte((message.ALT >> (8 * 3)) & 0xff)

				// TODO: sats
				payload[13] = 0

				for i := 0; i < len(payload); i++ {
					crc ^= payload[i]
				}

				b := []byte{0x24, 0x54, 0x47}
				b = append(b, payload...)
				b = append(b, crc)
				_, err := port.Write(b)
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}

			}
		})
	}

	robot := gobot.NewRobot("mavBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{iris},
		work,
	)

	robot.Start()
}
