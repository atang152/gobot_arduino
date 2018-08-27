package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func scale(v int, inMin int, inMax int, outMin int, outMax int) int {
	// Function to scale values from one range to another. Similar to map function used for Arduio
	return (v-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func elapsed(fn func(), sec time.Duration) {
	// Executes a function for a defined amount of seconds/duration. e.g. print hello world for 5 seconds

	start := time.Now()
	stop := time.After(sec * time.Second)
	for {
		select {
		case <-stop:
			fmt.Printf("EXIT: %v seconds\n", time.Since(start))
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Println(time.Since(start))
			fn()
		}
	}
}

func main() {

	sensorLow := 1023
	sensorHigh := 0

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	photoSensor := aio.NewAnalogSensorDriver(firmataAdaptor, "0") // Phototransistor sensor sensitive to light
	piezo := gpio.NewBuzzerDriver(firmataAdaptor, "8")            // Piezo sensor that converts electric to vibration and vice versa

	work := func() {

		// Calibration of Photosensor for 5 seconds
		start := time.Now()
		elapsed := time.Since(start)

		for elapsed < (time.Second * 5) {

			// Read phototransistor's raw value
			pSV, err := photoSensor.Read()
			if err != nil {
				fmt.Println(err)
			}

			if pSV > sensorHigh {
				sensorHigh = pSV
			}

			if pSV < sensorLow {
				sensorLow = pSV
			}

			fmt.Println("Sensor Low:", sensorLow, "\tSensor High:", sensorHigh, "\tSensor Value:", pSV)

			// Update elapse time
			elapsed = time.Since(start)
			time.Sleep(time.Second * 1)
		}
		fmt.Printf("Exit Calibration %v\n", elapsed)

		// Adjust pitch
		gobot.Every(1*time.Second, func() {

			pSV, err := photoSensor.Read()
			if err != nil {
				fmt.Println(err)
			}

			pitch := scale(pSV, sensorLow, sensorHigh, 50, 4000)
			fmt.Println("Sensor Low:", sensorLow, "\tSensor High:", sensorHigh, "\tSensor Value:", pSV, "\tPitch:", pitch)

			piezo.Tone(float64(pitch), 3)

			// JUST FOR FUN~~~
			// type note struct {
			// 	tone     float64
			// 	duration float64
			// }

			// song := []note{
			// 	{gpio.C4, gpio.Quarter},
			// 	{gpio.C4, gpio.Quarter},
			// 	{gpio.G4, gpio.Quarter},
			// 	{gpio.G4, gpio.Quarter},
			// 	{gpio.A4, gpio.Quarter},
			// 	{gpio.A4, gpio.Quarter},
			// 	{gpio.G4, gpio.Half},
			// 	{gpio.F4, gpio.Quarter},
			// 	{gpio.F4, gpio.Quarter},
			// 	{gpio.E4, gpio.Quarter},
			// 	{gpio.E4, gpio.Quarter},
			// 	{gpio.D4, gpio.Quarter},
			// 	{gpio.D4, gpio.Quarter},
			// 	{gpio.C4, gpio.Half},
			// }

			// for _, val := range song {
			// 	piezo.Tone(val.tone, val.duration)
			// 	time.Sleep(10 * time.Millisecond)
			// }

		})

	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{photoSensor, piezo},
		work,
	)
	robot.Start()
}

// ******************************************************************************************************************

// // Calibration of Photosensor for 5 seconds
// start := time.Now()
// stop := time.After(5 * time.Second)

// for {
// 	select {
// 	case <-stop:
// 		fmt.Printf("EXIT: %v seconds\n", time.Since(start))
// 		return
// 	default:
// 		// Read phototransistor's raw value
// 		time.Sleep(1 * time.Second)
// 		pSV, err := photoSensor.Read()
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		pSV, sensorHigh, sensorLow := calibrate(pSV, sensorHigh, sensorLow)
// 		fmt.Println("Sensor Low:", sensorLow, "\tSensor High:", sensorHigh, "\tSensor Value:", pSV)
// 	}
// }

// ******************************************************************************************************************

// func main() {
// 	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
// 	buzzer := gpio.NewBuzzerDriver(firmataAdaptor, "8")

// 	work := func() {
// 		type note struct {
// 			tone     float64
// 			duration float64
// 		}

// 		song := []note{
// 			{gpio.C4, gpio.Quarter},
// 			{gpio.C4, gpio.Quarter},
// 			{gpio.G4, gpio.Quarter},
// 			{gpio.G4, gpio.Quarter},
// 			{gpio.A4, gpio.Quarter},
// 			{gpio.A4, gpio.Quarter},
// 			{gpio.G4, gpio.Half},
// 			{gpio.F4, gpio.Quarter},
// 			{gpio.F4, gpio.Quarter},
// 			{gpio.E4, gpio.Quarter},
// 			{gpio.E4, gpio.Quarter},
// 			{gpio.D4, gpio.Quarter},
// 			{gpio.D4, gpio.Quarter},
// 			{gpio.C4, gpio.Half},
// 		}

// 		for _, val := range song {
// 			buzzer.Tone(val.tone, val.duration)
// 			time.Sleep(10 * time.Millisecond)
// 		}
// 	}

// 	robot := gobot.NewRobot("bot",
// 		[]gobot.Connection{firmataAdaptor},
// 		[]gobot.Device{buzzer},
// 		work,
// 	)

// 	robot.Start()
// }
