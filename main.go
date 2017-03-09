package main

//legg filene i GOPATH/src (finn ved å skrive go env i terminal.)

import (
	"fmt"
	"heisprosjekt/Events/ExternalOrders"
	"heisprosjekt/Network"
	. "heisprosjekt/Network/network/peers"
	"heisprosjekt/Timer"
	. "heisprosjekt/driver"
	"heisprosjekt/elevatorController/fsm"
	"heisprosjekt/elevatorController/runElevator"
	/*"time" */)

func main() {
	elevatorData := runElevator.InitializeElevator()
	fmt.Println(ExternalOrders.Elevators)
	ExternalOrders.UpdateElevatorData(elevatorData)

	updateElevatorRxCh := make(chan ElevatorData, 10)
	updateElevatorTxCh := make(chan ElevatorData, 10)

	startTimer := make(chan TimerType, 10)
	timeOut := make(chan TimerType, 10)

	newOrderTxCh := make(chan ElevatorOrder, 10)
	newOrderRxCh := make(chan ElevatorOrder, 10)

	peerUpdateCh := make(chan PeerUpdate, 10)
	peerTxEnableCh := make(chan bool)

	arriveAtFloorCh := make(chan int)
	externalButtonCh := make(chan ElevatorOrder, 10)
	internalButtonCh := make(chan int, 10)

	go Network.RunNetwork(elevatorData, updateElevatorTxCh, updateElevatorRxCh, newOrderTxCh, newOrderRxCh, peerUpdateCh, peerTxEnableCh)

	go runElevator.ReadFloorSensors(arriveAtFloorCh)
	go runElevator.ReadButtonSensors(externalButtonCh, internalButtonCh)

	go Timer.RunTimer(timeOut, startTimer)

	for {
		select {

		case msg1 := <-arriveAtFloorCh:
			//fsmArriveAtFloor(msg)

			elevatorData = fsm.ArriveAtFloor(elevatorData, msg1, startTimer)

		case msg2 := <-externalButtonCh:
			//elevatorData = fsmExternalButtonPressed(elevatorData, msg)
			elevatorData = fsm.ExternalButtonPressed(elevatorData, msg2, newOrderTxCh, updateElevatorTxCh, startTimer)

		case msg3 := <-internalButtonCh:
			elevatorData = fsm.InternalButtonPressed(elevatorData, msg3, updateElevatorTxCh, startTimer)

		case msg4 := <-updateElevatorRxCh:
			elevatorData = fsm.ElevatorUpdateReceived(msg4, elevatorData)

			for i := 0; i < 3; i++ {
			}

		case msg5 := <-newOrderRxCh:
			elevatorData = fsm.NewOrderReceived(msg5, elevatorData, updateElevatorTxCh)

			//elevatorData = OrderReceivedOrder(elevatorData, msg)
		case msg6 := <-peerUpdateCh:
			fmt.Println("PeerUpdate: ", msg6)
			go ExternalOrders.OnlineElevatorsUpdate(msg6, updateElevatorTxCh, newOrderTxCh)

		case timeout := <-timeOut:
			elevatorData = fsm.TimeOut(elevatorData, timeout)

		}

	}

}

/*
func main1() {
	NetworkTest()

}

func NetworkTest() {

	elevatorData := InitializeElevator()

	updateElevatorRxCh := make(chan ElevatorData, 50)
	updateElevatorTxCh := make(chan ElevatorData, 50)

	startTimer := make(chan TimerType, 10)
	timeOut := make(chan TimerType, 10)

	newOrderTxCh := make(chan ElevatorOrder, 50)
	newOrderRxCh := make(chan ElevatorOrder, 50)

	peerUpdateCh := make(chan PeerUpdate, 50)
	peerTxEnableCh := make(chan bool)

	arriveAtFloorCh := make(chan int)
	externalButtonCh := make(chan ElevatorOrder, 50)
	internalButtonCh := make(chan int, 50)

	go RunNetwork(elevatorData, updateElevatorTxCh, updateElevatorRxCh, newOrderTxCh, newOrderRxCh, peerUpdateCh, peerTxEnableCh)

	go ReadAllSensors2(arriveAtFloorCh, externalButtonCh, internalButtonCh)

	for {
		select {

		case msg1 := <-arriveAtFloorCh:
			//fsmArriveAtFloor(msg)

			elevatorData = FsmArriveAtFloor(elevatorData, msg1, startTimer)

		case msg2 := <-externalButtonCh:
			var testOrder ElevatorOrder
			testOrder.Floor = 0
			testOrder.Direction = 0
			testOrder.ElevatorID = "test"

			newOrderTxCh <- testOrder
			fmt.Println(msg2)

		case msg3 := <-internalButtonCh:
			updateElevatorTxCh <- elevatorData
			fmt.Println(msg3)
		case msg4 := <-updateElevatorRxCh:
			fmt.Println("Elevator Update received")
			fmt.Println(msg4)
			//elevatorData = OrderReceivedUpdate(elevatorData, msg)

		case msg5 := <-newOrderRxCh:
			fmt.Println("New order received")
			fmt.Println(msg5)
			//elevatorData = OrderReceivedOrder(elevatorData, msg)
		case msg6 := <-peerUpdateCh:
			fmt.Println(msg6)
			//elevatorData = PeerUpdate(elevatorData, msg)

		}
	}

}

*/
