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

	const sensorPin = "0"

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	sensor := aio.NewAnalogSensorDriver(firmataAdaptor, sensorPin)
	ledG := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		gobot.Every(1*time.Second, func() {

			sensorVal, err := sensor.Read()

			if err != nil {
				fmt.Println(err)
			}

			volt := (float64(sensorVal) * 5.0) / 1024.0
			temp := (volt - 0.5) * 100

			if temp > 32 {
				ledG.On()
			} else {
				fmt.Println("Turn off Led")
				ledG.Off()
			}

			fmt.Println("Voltage:", volt, "sensorVal:", sensorVal, "Temperature:", temp)

		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()

}
