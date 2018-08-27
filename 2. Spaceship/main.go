package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	button := gpio.NewButtonDriver(firmataAdaptor, "2")
	ledR5 := gpio.NewLedDriver(firmataAdaptor, "5")
	ledR4 := gpio.NewLedDriver(firmataAdaptor, "4")
	ledG3 := gpio.NewLedDriver(firmataAdaptor, "3")

	work := func() {

		ledG3.On()
		fmt.Println("Green Led is on")

		button.On(gpio.ButtonPush, func(data interface{}) {

			ledG3.Off()
			fmt.Println("Button is pressed. Green Led is OFF")

			for button.Active == true {
				ledR4.Toggle()
				ledR5.Toggle()
				fmt.Println("Red Led is ON", ledR5.State())
				time.Sleep(2 * time.Second)
				fmt.Println("Red Led is OFF", ledR5.State())
			}

		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Button is Release. Green Led is ", ledG3.State())
			ledG3.On()
			ledR4.Off()
			ledR5.Off()
		})

	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{button, ledR5, ledR4, ledG3},
		work,
	)

	robot.Start()
}
