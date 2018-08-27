package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	gLed := gpio.NewLedDriver(firmataAdaptor, "9")
	rLed := gpio.NewLedDriver(firmataAdaptor, "10")
	bLed := gpio.NewLedDriver(firmataAdaptor, "11")

	rSensor := aio.NewAnalogSensorDriver(firmataAdaptor, "0")
	gSensor := aio.NewAnalogSensorDriver(firmataAdaptor, "1")
	bSensor := aio.NewAnalogSensorDriver(firmataAdaptor, "2")

	work := func() {
		gobot.Every(1*time.Second, func() {

			// Read value (0-1023) off Sensor and convert it to duty cycle (0-255)
			rSV, err := rSensor.Read()
			if err != nil {
				fmt.Println(err)
			}
			rV := byte(rSV / 4)

			// Read value (0-1023) off Sensor and convert it to duty cycle (0-255)
			gSV, err := gSensor.Read()
			if err != nil {
				fmt.Println(err)
			}
			gV := byte(gSV / 4)

			// Read value (0-1023) off Sensor and convert it to duty cycle (0-255)
			bSV, err := bSensor.Read()
			if err != nil {
				fmt.Println(err)
			}
			bV := byte(bSV / 4)

			// fmt.Println("Red Sensor Value:", rSV, "\t Green Sensor Value:", gSV, "\t Blue Sensor Value:", bSV)
			fmt.Println("Red Value:", rV, "\t Green Value:", gV, "\t Blue Value:", bV)

			// Write duty cycle value (0-255) as byte into Brightness method that takes a pointer to LedDriver
			rLed.Brightness(rV)
			gLed.Brightness(gV)
			bLed.Brightness(bV)
		})

	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{gLed, rLed, bLed, rSensor, gSensor, bSensor},
		work,
	)

	robot.Start()
}
