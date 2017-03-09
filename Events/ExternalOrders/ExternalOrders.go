package ExternalOrders

import (
	"fmt"

	. "heisprosjekt/Network/network/peers"
	. "heisprosjekt/driver"
)

const MASTER = false

var ThisElevatorID = ""

var NETWORK_DOWN = false //hvordan er det bra å bruke denne variablen

var ExternalOrderLights = make([]int, 0) //skal den være her eller i intern-delen?

var Elevators = make([]ElevatorData, N_ELEVATORS)

func OnlineElevatorsUpdate(onlineElevatorList PeerUpdate, updateElevatorTxCh chan ElevatorData, newOrderCh chan ElevatorOrder) {

	//Setting all lost elevators to uninitiated, is probably unessecary.Unless two elevators fail at the same instant
	for i := 0; i < len(onlineElevatorList.Lost); i++ {
		Elevators[FindElevatorIndex(onlineElevatorList.Lost[i])].Initiated = false
	}

	if onlineElevatorList.New == "" {
		RedestributeExternalOrders(Elevators[FindElevatorIndex(onlineElevatorList.Lost[len(onlineElevatorList.Lost)-1])], newOrderCh, updateElevatorTxCh)
	}

	if onlineElevatorList.New != "" {
		if FindElevatorIndex(onlineElevatorList.New) == -1 {
			Elevators[FindElevatorIndex("")].ID = onlineElevatorList.New
			Elevators[FindElevatorIndex(onlineElevatorList.New)].Initiated = true
		} else {

			Elevators[FindElevatorIndex(onlineElevatorList.New)].Initiated = true
			//If this elevator already has data stored we want to push those data back
			if HasUnresolvedInternalOrders(Elevators[FindElevatorIndex(onlineElevatorList.New)]) == true {
				fmt.Println("her her", Elevators[FindElevatorIndex(onlineElevatorList.New)])
				Elevators[FindElevatorIndex(onlineElevatorList.New)].ForceUpdate = true //To indicate that this data packet should be forced
				updateElevatorTxCh <- Elevators[FindElevatorIndex(onlineElevatorList.New)]
			}
		}
	}

	if len(onlineElevatorList.Peers) == 1 && len(onlineElevatorList.Lost) > 1 {
		//Lost connection to all other Elevators
		//Initialize reboot
	}

}

/*func OnlineElevatorsUpdate2(onlineElevatorList PeerUpdate) {

  //Dersom det ikke er nye har vi nødvendigvis mistet en heis
  if onlineElevatorList.New == "" {

    for i := 0; i < N_ELEVATORS; i++ {
      if onlineElevatorList.Lost[len(onlineElevatorList.Lost)-1] == Elevators[i].ID {
        Elevators[i].Initiated = false
      }
    }
  } else if onlineElevatorList.New == onlineElevatorList.Peers[0] {
    Elevators[0].ID = onlineElevatorList.New
    ThisElevatorID = onlineElevatorList.New
  } else {

    for i := 0; i < N_ELEVATORS; i++ {
      if Elevators[i].ID == onlineElevatorList.New {
        Elevators[i].Initiated = true
      }
    }
    for i := 0; i < N_ELEVATORS; i++ {
      if Elevators[i].ID == "" {
        Elevators[i].ID = onlineElevatorList.New
        Elevators[i].Initiated = true
      } else {
        fmt.Print("Noe er galt i OnlineElevatorsUpdate")
      }
    }
  }

}*/

func CalculateSingleElevatorCost(elevator ElevatorData, order ElevatorOrder) int {
	if (int(elevator.Direction) == -1 && int(order.Direction) == 1) || (int(elevator.Direction) == 1 && int(order.Direction) == 0) { //her blir det krøll
		switch elevator.Direction {
		case DirnUp:
			if order.Floor > elevator.Floor {
				return order.Floor - elevator.Floor
			} else {
				return (elevator.Floor-1)*2 + (elevator.Floor - order.Floor)
			}
		case DirnDown:
			if order.Floor < elevator.Floor {
				return elevator.Floor - order.Floor
			} else {
				return (elevator.Floor-1)*2 + (order.Floor - elevator.Floor)
			}
		}
	} else {
		switch elevator.Direction {
		case DirnUp:
			return 2*N_FLOORS - elevator.Floor - order.Floor
		case DirnDown:
			return (elevator.Floor - 1) + (order.Floor - 1)
		}
	}
	return -1
}

func FindBestElevator(order ElevatorOrder) string {
	var minCost = 100000
	var ID string

	//var elevatorNumber = -1 //kanksje fint å bruke ID her?
	for i := 0; i < N_ELEVATORS; i++ {
		if Elevators[i].Initiated {
			var thisCost = CalculateSingleElevatorCost(Elevators[i], order)
			if thisCost < minCost {
				minCost = thisCost
				ID = Elevators[i].ID
			}
		}
	}
	return ID //kan bare bruke ID-en til ordren
}

//Dette må også ordnes 23feb
/*func PlaceExternalOrder2(elevatorData ElevatorData, order ElevatorOrder) ElevatorData {
  elevatorData.Orders[order.Floor][order.Direction] = 1
  return elevatorData
}*/

func PlaceInternalOrder(elevatorData ElevatorData, floor int, updateElevatorTxCh chan ElevatorData) ElevatorData {

	elevatorData.Orders[floor][ButtonType(2)] = 1
	updateElevatorTxCh <- elevatorData

	return elevatorData
}

func PlaceExternalOrder(elevatorData ElevatorData, order ElevatorOrder, newOrderTxCh chan ElevatorOrder, updateElevatorTxCh chan ElevatorData) ElevatorData {
	order.ElevatorID = FindBestElevator(order)
	if order.ElevatorID == elevatorData.ID {
		//Oppdaterer egne ordreliste
		elevatorData.Orders[order.Floor][order.Direction] = 1

		//Sender oppdatert informasjon på nettverket
		updateElevatorTxCh <- elevatorData

	} else {

		newOrderTxCh <- order

	}
	return elevatorData

}

/*func SuccessfulPlacementConfirmation(elevatorNumber int, order ElevatorOrder) bool {
  if Elevators[elevatorNumber].Orders[order.Floor][order.Direction] == 1 {
    return true
  }
  return false
}*/

//må lage noe som merker at en heis har falt ut. -- lages i nettverk

//tror det er best om bare en av heisene omfordeler ordre

func RedestributeExternalOrders(lostElevator ElevatorData, newOrderCh chan ElevatorOrder, updateElevatorDataCh chan ElevatorData) {

	if AmIMaster() == true {
		fmt.Println("redistributer")
		fmt.Println(lostElevator)
		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < 2; j++ {
				if lostElevator.Orders[i][j] == 1 {
					Elevators[FindElevatorIndex(lostElevator.ID)].Orders[i][j] = 0
					newOrder := ElevatorOrder{i, j, ""}
					newOrder.ElevatorID = FindBestElevator(newOrder)

					newOrderCh <- newOrder
				}
			}
		}

		updateElevatorDataCh <- Elevators[FindElevatorIndex(lostElevator.ID)]
	}
}

func DenyNewExternalOrders(elevatorData ElevatorData) { //on network fall out
	for i := 0; i < N_FLOORS; i++ {
		elevatorData.Orders[i][0] = 0
		elevatorData.Orders[i][1] = 0
	}
	NETWORK_DOWN = true
}

func FindElevatorIndex(elevatorID string) int {
	index := -1
	for i := 0; i < N_ELEVATORS; i++ {
		if Elevators[i].ID == elevatorID {
			return i
		}
	}
	return index
}

func UpdateElevatorData(elevatorData ElevatorData) {
	if FindElevatorIndex(elevatorData.ID) == -1 {
		panic("Trying to access elevator that is not in the Elevators list doesnt exist")
	}
	Elevators[FindElevatorIndex(elevatorData.ID)] = elevatorData
}

func InitializeElevatorList(elevatorData ElevatorData) {
	Elevators[0] = elevatorData
	for i := 1; i < N_ELEVATORS; i++ {
		Elevators[i].ID = ""
		for j := 0; j < N_FLOORS; j++ {
			for k := 0; k < N_BUTTONS; k++ {
				Elevators[i].Orders[j][k] = 0
			}
		}
		Elevators[i].Initiated = false
	}
}

func AllExternalOrders() [N_FLOORS][N_BUTTONS]int {
	var allExternalOrders [N_FLOORS][N_BUTTONS]int
	for i := 0; i < N_ELEVATORS; i++ {
		for j := 0; j < N_FLOORS; j++ {
			for k := 0; k < 2; k++ {
				if Elevators[i].Orders[j][k] == 1 {
					allExternalOrders[j][k] = Elevators[i].Orders[j][k]
				}
			}
		}
	}
	return allExternalOrders
}

func HasUnresolvedInternalOrders(elevatorData ElevatorData) bool {
	for i := 0; i < N_FLOORS; i++ {
		if elevatorData.Orders[i][2] == 1 {
			return true
		}
	}

	return false
}

//------------------------------------------------------------------------------
//Lagt til av Erling

//In case of updated list of connected elevators
/*func EventUpdatedPeers(updatedConnectionData PeerUpdate) {

	//Either we have fewer connections or more connections. Either way
	//we want to update our

}*/

/*func HasOrdersAtCurrentFloor(elevatorData ElevatorData) bool {
  return (elevatorData.Orders[elevatorData.Floor][0] == 1 || elevatorData.Orders[elevatorData.Floor][1] == 1 || elevatorData.Orders[elevatorData.Floor][2] == 1) && elevatorData.AtFloor == true
}*/

func AmIMaster() bool {
	for i := 1; i < len(Elevators); i++ {
		if Elevators[i].Initiated == true && Elevators[i].ID > Elevators[0].ID {
			return false
		}

	}

	return true

}
