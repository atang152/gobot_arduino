package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// long map(long x, long in_min, long in_max, long out_min, long out_max)
// {
//   return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min;
// }

func scale(v int, inMin int, inMax int, outMin int, outMax int) int {
	// function to scale values from one range to another. Similar to map function used for Arduio
	return (v-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	potM := aio.NewAnalogSensorDriver(firmataAdaptor, "0")
	servo := gpio.NewServoDriver(firmataAdaptor, "9")

	work := func() {
		// Set servo to it's minimum position which is 0
		servo.Min()

		gobot.Every(1*time.Second, func() {

			// Potentiometer Values return a value from 0 to 1024
			// Therefore, we would have to convert it to an angle
			// degree from 0 to 180
			potV, err := potM.Read()
			if err != nil {
				fmt.Println(err)
			}
			angle := scale(potV, 0, 1023, 0, 179)
			fmt.Println("Potentiometer Raw Value:", potV, "\t Angle:", angle)

			// To control the servo rotor which takes a angle from 0 to 180
			i := uint8(angle)
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}
