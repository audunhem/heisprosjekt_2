package Timer

import (
  . "heisprosjekt/driver"
  "time"
)

const openDoorTime = 3
const reachingFloorTime = 3

func RunTimer(timeOut chan TimerType, startTimer chan TimerType) {

  timer := time.NewTimer(0)
  timer.Stop()

  for {

    select {

    case timerType := <-startTimer:

      //If we want to start the timer for open doors
      if timerType == TimeToOpenDoors {
        go SetOpenDoorTimer(timeOut)

        //If we want to start the timer to see how long
        //to reach floor
      } else if timerType == TimeToReachFloor {
        timer.Reset(reachingFloorTime * time.Second)

        //If we have reached the floor and want to stop the timer
      } else if timerType == TimeFloorReached {
        timer.Stop()
      }

    case <-timer.C:
      timeOut <- TimerType(0)
    }

  }

}

func SetOpenDoorTimer(timeOut chan TimerType) {

  time.Sleep(openDoorTime * time.Second)

  timeout := TimerType(1)

  timeOut <- timeout

}
