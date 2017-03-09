package fsm

import (
	"fmt"

	. "heisprosjekt/Events/ElevatorNextFloorControl"
	. "heisprosjekt/Events/ExternalOrders"
	. "heisprosjekt/driver"
)

//Funksjoner som skal legges til ORDER MODULEN
//------------------------------------------------------------------------------------------------------------

func ArriveAtFloor(elevatorData ElevatorData, floor int, startTimer chan TimerType) ElevatorData {

	if floor == -1 {
		elevatorData.AtFloor = false
		//Vi har forlatt en etasje, starter timeren
		startTimer <- TimerType(TimeToReachFloor)
		//elevatorData = RemoveCompletedOrders(elevatorData)
		//SetAllLights(elevatorData, AllExternalOrders())
	} else {

		startTimer <- TimerType(TimeFloorReached)
		SetFloorIndicator(floor)
		elevatorData.AtFloor = true
		elevatorData.Floor = floor
		PrintOrderList(elevatorData)
		if CheckIfShouldStop(elevatorData) == true {
			elevatorData.Status = StatusDoorOpen
			SetMotorDirection(DirnStop)
			SetDoorOpenLamp(1)
			//elevatorData = OpenDoors(elevatorData)
			startTimer <- TimerType(TimeToOpenDoors)

			//elevatorData = RemoveCompletedOrders(elevatorData)
			PrintOrderList(elevatorData)
		}

		UpdateElevatorData(elevatorData)
		//SetAllLights(elevatorData, AllExternalOrders())

	}
	return elevatorData
}

func ExternalButtonPressed(elevatorData ElevatorData, order ElevatorOrder, newOrderTxCh chan ElevatorOrder, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) ElevatorData {

	elevatorData = PlaceExternalOrder(elevatorData, order, newOrderTxCh, updateElevatorTxCh)

	if elevatorData.Direction == DirnStop && !DoorOpenLampOn() {
		elevatorData = OrderSetNextDirection(elevatorData)
		SetMotorDirection(elevatorData.Direction)
	}

	if !NoOrdersAtCurrentFloor(elevatorData) {
		elevatorData = ArriveAtFloor(elevatorData, elevatorData.Floor, startTimer)
	}

	//	if elevatorData.Status == StatusIdle && elevatorData.Statys {}

	//	if elevatorData.Status == StatusIdle && elevatorData.Statys {}

	UpdateElevatorData(elevatorData)
	SetAllLights(elevatorData, AllExternalOrders())

	return elevatorData

}

/*func FsmStopAtFloor() {
	SetMotorDirection(DirnStop)
	SetDoorOpenLamp(1)
	time.Sleep(500 * time.Millisecond)
	SetDoorOpenLamp(0)
}*/

/*func OpenDoors(elevatorData ElevatorData) ElevatorData {
	elevatorData.Status = StatusDoorOpen
	SetMotorDirection(DirnStop)
	SetDoorOpenLamp(1)
	return elevatorData
}*/

func LeaveFloor(elevatorData ElevatorData) ElevatorData {
	elevatorData = OrderSetNextDirection(elevatorData)
	elevatorData = RemoveCompletedOrders(elevatorData)
	SetDoorOpenLamp(0)
	SetAllLights(elevatorData, AllExternalOrders())
	elevatorData = OrderSetNextDirection(elevatorData)
	PrintOrderList(elevatorData)
	SetMotorDirection(elevatorData.Direction)
	return elevatorData
}

func PrintOrderList(elevatorStruct ElevatorData) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < N_BUTTONS; j++ {
			fmt.Printf("%d", elevatorStruct.Orders[i][j])
		}
		fmt.Printf("\n")

	}

	fmt.Printf("--------------------------------------------")
	fmt.Printf("\n")

}

func InternalButtonPressed(elevatorStruct ElevatorData, floor int, updateElevatorTxCh chan ElevatorData, startTimer chan TimerType) ElevatorData {

	elevatorData := PlaceInternalOrder(elevatorStruct, floor, updateElevatorTxCh)

	if elevatorData.Direction == DirnStop {

		if elevatorData.Floor == floor {
			elevatorData = ArriveAtFloor(elevatorData, floor, startTimer)
		} else {
			//fmt.Println("Internal")
			elevatorData = OrderSetNextDirection(elevatorData)
			SetMotorDirection(elevatorData.Direction)
			//elevatorData = RemoveCompletedOrders(elevatorData)
			//SetAllLights(elevatorData, AllExternalOrders())
		}

	}
	//UpdateElevatorData(elevatorData)
	SetAllLights(elevatorData, AllExternalOrders())

	return elevatorData

}

func ElevatorUpdateReceived(elevatorDataRx ElevatorData, elevatorData ElevatorData) ElevatorData {
	if elevatorDataRx.ID == elevatorData.ID && elevatorDataRx.ForceUpdate == true {
		elevatorData.Orders = elevatorDataRx.Orders
		elevatorData.ForceUpdate = false
		fmt.Println("Has received unresolved internal orders: ")
		PrintOrderList(elevatorData)
		elevatorData = OrderSetNextDirection(elevatorData)
		SetMotorDirection(elevatorData.Direction)
		//UpdateElevatorData(elevatorData)

	} else {
		UpdateElevatorData(elevatorDataRx)
	}

	SetAllLights(elevatorData, AllExternalOrders())

	return elevatorData
}

func NewOrderReceived(order ElevatorOrder, elevatorData ElevatorData, elevatorUpdateTxCh chan ElevatorData) ElevatorData {

	if order.ElevatorID == elevatorData.ID {
		//The order belongs to this elevator
		elevatorData.Orders[order.Floor][order.Direction] = 1

		//Broadcast the updated elevator struct
		elevatorUpdateTxCh <- elevatorData

		if elevatorData.Direction == DirnStop {
			elevatorData = OrderSetNextDirection(elevatorData)
			SetMotorDirection(elevatorData.Direction)
		}

		UpdateElevatorData(elevatorData)
		SetAllLights(elevatorData, AllExternalOrders())

	}

	return elevatorData

}

func SetAllLights(elevatorData ElevatorData, allExternalOrders [N_FLOORS][N_BUTTONS]int) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < N_BUTTONS-1; j++ {
			SetButtonLamp(ButtonType(j), i, allExternalOrders[i][j])
			SetButtonLamp(ButtonType(2), i, elevatorData.Orders[i][2])
		}
	}
}

func TimeOut(elevatorData ElevatorData, timeout TimerType) ElevatorData {

	if timeout == TimeToOpenDoors {
		elevatorData = LeaveFloor(elevatorData)
	} else if timeout == TimeToReachFloor {
		panic("Cant reach floor, sensor/engine error!")
	}

	return elevatorData
}

/*

func goDown() {}

func goUp() {}

func openDoors() {}

func stop() {}

func readAllSensors() {}

//DETTE ER VAR FØRSTE UTKASTET PRØVER PÅ NYTT
/*


func GoToFloor(floor int, updatedData *ElevatorData) bool {
	fmt.Println("GoToFloor")
	if (floor - (*updatedData).Floor ) > 0 {
		SetMotorDirection(DirnUp)
		fmt.Println("MotorDirectionUp")
	} else if (floor - (*updatedData).Floor ) < 0 {
		SetMotorDirection(DirnDown)
		fmt.Println("MotorDirectionDown")
	}
	for (*updatedData).Floor != floor {
	}
	fmt.Println("GoToFloorEnd")
	SetMotorDirection(DirnStop)
	OpenDoors()
	fmt.Println("GoToFloorEnd")

	return true
}

func OpenDoors() {
	if GetMotorDirection() != DirnStop {
		fmt.Println("Heisen har ikke stoppet")
	}
	SetDoorOpenLamp(1)
	time.Sleep(3 * time.Second)
	SetDoorOpenLamp(0)
}

func GoUp...() {

}


/*
import (
	."./../driver"
	"fmt"
	"time"

)

func GoToFloor(floor int){
	//Takes the elevator to _floor_, opens the door for 3 secs and closes it. Returns 1 on success



	if (GetFloorSensorSignal() == -1) {
		SetMotorDirection(DirnUp)

			for (GetFloorSensorSignal() == -1) {}
			SetMotorDirection(DirnStop)
	}




	if (floor - GetFloorSensorSignal()) > 0 {
		SetMotorDirection(DirnUp)
	} else if (floor - GetFloorSensorSignal()) < 0 {
		SetMotorDirection(DirnDown)
	}

	fmt.Println(GetMotorDirection())
	for (GetFloorSensorSignal() != floor) {
	}




	SetMotorDirection(DirnStop)
	OpenDoors()



*/

//else if (floor == GetFloorSensorSignal())

/*
}


func OpenDoors() {
	if GetMotorDirection() != DirnStop {
		fmt.Println("gjør noe")
	}
	SetDoorOpenLamp(1)
	time.Sleep(3*time.Second)
	SetDoorOpenLamp(0)
}*/
