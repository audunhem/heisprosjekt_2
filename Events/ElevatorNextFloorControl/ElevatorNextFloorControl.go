package ElevatorNextFloorControl

import (
  . "heisprosjekt/Events/ExternalOrders"
  . "heisprosjekt/driver"
  //"fmt"
)

var thisElevator = Elevators[0]

//trenger noe ala thisElevator

/*func CheckIfShouldStop(elevatorData ElevatorData) bool {
  switch {
  case elevatorData.Direction == DirnUp:
    if elevatorData.Orders[elevatorData.Floor][ButtonCallUp] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
      return true
    } else if elevatorData.Floor == N_FLOORS-1 {
      return true

    } else {
      for i := elevatorData.Floor + 1; i < N_FLOORS; i++ {
        if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
          return false
        }
      }
      return true
    }
    return false
  case elevatorData.Direction == DirnDown:
    if elevatorData.Orders[elevatorData.Floor][ButtonCallDown] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
      //elevatorData.Orders[elevatorData.Floor][ButtonCallUp] = false
      //elevatorData.Orders[elevatorData.Floor][ButtonInternal] = false
      //mulig dette kan føre til at ordre forsvinner, og kanskje bedre med en egen funksjon for funksjonaliteten
      return true
    } else if elevatorData.Floor == 0 {
      return true
    } else {
      for i := 0; i < elevatorData.Floor; i++ {
        if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
          return false
        }
      }
      return true
    }
    return false
  }
  return false
}*/

//må kalles etter "dørene lukkes" og neste retning er satt
func RemoveCompletedOrders(elevatorData ElevatorData) ElevatorData {
  switch elevatorData.Direction {

  case DirnUp:

    elevatorData.Orders[elevatorData.Floor][ButtonCallUp] = 0
    elevatorData.Orders[elevatorData.Floor][ButtonInternal] = 0

    if NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Orders[elevatorData.Floor][ButtonCallDown] = 0 //hvis de som skal opp ikke trykker videre, slettes denne, og det er litt uheldig
    }

  case DirnDown:

    elevatorData.Orders[elevatorData.Floor][ButtonCallDown] = 0
    elevatorData.Orders[elevatorData.Floor][ButtonInternal] = 0

    if NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Orders[elevatorData.Floor][ButtonCallUp] = 0
    }

  case DirnStop:
    if !NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Orders[elevatorData.Floor][ButtonCallDown] = 0
    } else if !NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Orders[elevatorData.Floor][ButtonCallUp] = 0
    }
    elevatorData.Orders[elevatorData.Floor][ButtonInternal] = 0

  }
  UpdateElevatorData(elevatorData)
  return elevatorData
}

func OrderSetNextDirection(elevatorData ElevatorData) ElevatorData {

  if NoOrdersAboveCurrentFloor(elevatorData) && NoOrdersAtCurrentFloor(elevatorData) && NoOrdersBelowCurrentFloor(elevatorData) {
    elevatorData.Direction = DirnStop
    goto end
  }

  switch elevatorData.Direction {

  case DirnUp:

    if NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnDown
    } else {
      elevatorData.Direction = DirnUp
    }

  case DirnDown:

    if NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnUp
    } else {
      elevatorData.Direction = DirnDown
    }

  case DirnStop:

    if !NoOrdersAtCurrentFloor(elevatorData) {
      //to stop the elevator from picking up all orders at a floor
      if !NoOrdersBelowCurrentFloor(elevatorData) {
        elevatorData.Direction = DirnDown
      } else if !NoOrdersAboveCurrentFloor(elevatorData) {
        elevatorData.Direction = DirnUp
      }
      elevatorData.Direction = DirnStop
    }

    if NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnUp
    } else {
      elevatorData.Direction = DirnDown
    }

  }

end:
  UpdateElevatorData(elevatorData)
  return elevatorData
}

/*func OrderSetNextDirection(elevatorData ElevatorData) ElevatorData {

  check := 0

  if elevatorData.Status == StatusIdle {

    if !NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnUp
      SetMotorDirection(DirnUp)
      elevatorData.Status = StatusMoving
      check = 1

    } else if !NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnDown
      SetMotorDirection(DirnDown)
      elevatorData.Status = StatusMoving
      check = 1
    }

  } else if elevatorData.Direction == DirnUp {
    if !NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnUp
      SetMotorDirection(DirnUp)
      elevatorData.Status = StatusMoving
      check = 1

    } else if !NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnDown
      SetMotorDirection(DirnDown)
      elevatorData.Status = StatusMoving
      check = 1
    }

  } else if elevatorData.Direction == DirnDown {
    if !NoOrdersBelowCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnDown
      SetMotorDirection(DirnDown)
      elevatorData.Status = StatusMoving
      check = 1

    } else if !NoOrdersAboveCurrentFloor(elevatorData) {
      elevatorData.Direction = DirnUp
      SetMotorDirection(DirnUp)
      elevatorData.Status = StatusMoving
      check = 1

    }

  }

  if check == 0 {
    elevatorData.Direction = DirnStop
    elevatorData.Status = StatusIdle
  }
  return elevatorData
}

func OrderSetNextDirection2(elevatorStruct ElevatorData) ElevatorData {
  elevatorData := elevatorStruct
  check := 0

  if elevatorData.Status == StatusIdle {
    for i := 0; i < N_FLOORS; i++ {
      for j := 0; j < N_BUTTONS; j++ {
        if elevatorData.Orders[i][j] == 1 {
          if elevatorData.Floor < i {
            elevatorData.Direction = DirnUp
            SetMotorDirection(DirnUp)
            elevatorData.Status = StatusMoving
          } else if elevatorData.Floor > i {
            elevatorData.Direction = DirnDown
            SetMotorDirection(DirnDown)
            elevatorData.Status = StatusMoving
          } else if elevatorData.Floor == i {
            //elevatorData = OpenDoors()
          }
        }

      }
    }

  } else if elevatorData.Direction == DirnUp {

    for i := elevatorData.Floor; i < N_FLOORS; i++ {
      for j := 0; j < N_BUTTONS; j++ {
        if elevatorData.Orders[i][j] == 1 {
          SetMotorDirection(DirnUp)
          check = 1
        }
      }
    }

    if check == 0 {
      for i := 0; i < elevatorData.Floor; i++ {
        for j := 0; j < N_BUTTONS; j++ {
          if elevatorData.Orders[i][j] == 1 {
            SetMotorDirection(DirnDown)
            elevatorData.Direction = DirnDown
            check = 1
          }
        }
      }
    }

    if check == 0 {
      elevatorData.Status = StatusIdle
      elevatorData.Direction = DirnStop
    }

  } else if elevatorData.Direction == DirnDown {
    for i := 0; i < elevatorData.Floor; i++ {
      for j := 0; j < N_BUTTONS; j++ {
        if elevatorData.Orders[i][j] == 1 {
          SetMotorDirection(DirnDown)
          check = 1
        }
      }
    }

    if check == 0 {
      for i := elevatorData.Floor; i < N_FLOORS; i++ {
        for j := 0; j < N_BUTTONS; j++ {
          if elevatorData.Orders[i][j] == 1 {
            SetMotorDirection(DirnUp)
            elevatorData.Direction = DirnUp
            check = 1
          }
        }
      }
    }

    if check == 0 {
      elevatorData.Status = StatusIdle
      elevatorData.Direction = DirnStop
    }

  } else {
    elevatorData.Direction = DirnStop
  }

  return elevatorData
}*/

func CheckIfShouldStop(elevatorData ElevatorData) bool {
  switch elevatorData.Direction {

  case DirnUp:
    if elevatorData.Orders[elevatorData.Floor][ButtonCallUp] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
      return true
    } else if NoOrdersAboveCurrentFloor(elevatorData) {
      return true
    }

  case DirnDown:
    if elevatorData.Orders[elevatorData.Floor][ButtonCallDown] == 1 || elevatorData.Orders[elevatorData.Floor][ButtonInternal] == 1 {
      return true
    } else if NoOrdersBelowCurrentFloor(elevatorData) {
      return true
    }

  case DirnStop:
    return true

  }
  return false
}

func NoOrdersAboveCurrentFloor(elevatorData ElevatorData) bool {
  if elevatorData.Floor == N_FLOORS-1 {
    return true
  }
  for i := elevatorData.Floor + 1; i < N_FLOORS; i++ {
    if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
      return false
    }
  }
  return true
}

func NoOrdersBelowCurrentFloor(elevatorData ElevatorData) bool {
  if elevatorData.Floor == 0 {
    return true
  }
  for i := 0; i < elevatorData.Floor; i++ {
    if elevatorData.Orders[i][ButtonCallUp] != 0 || elevatorData.Orders[i][ButtonCallDown] != 0 || elevatorData.Orders[i][ButtonInternal] != 0 {
      return false
    }
  }
  return true
}

func NoOrdersAtCurrentFloor(elevatorData ElevatorData) bool {
  for i := 0; i < N_BUTTONS; i++ {
    if elevatorData.Orders[elevatorData.Floor][i] != 0 {
      return false
    }
  }
  return true
}
