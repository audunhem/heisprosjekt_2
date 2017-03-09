package main

import (
	. "./driver"
	"fmt"
)

func main() {
	InitElevator()
	fmt.Println("Press STOP button to stop elevator and exit program.")
	SetMotorDirection(DirnUp)
	for {
		GetButtonSignal(0, 2)
		if GetFloorSensorSignal() == N_FLOORS-1 {
			SetMotorDirection(DirnDown)

		} else if GetFloorSensorSignal() == 0 {
			SetMotorDirection(DirnUp)
		}

		// Stop elevator and exit program if the stop button is pressed
		if GetStopSignal() == 1 {
			SetMotorDirection(DirnStop)
		}
	}
}
