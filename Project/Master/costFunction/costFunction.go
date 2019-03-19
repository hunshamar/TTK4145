package costFunction

import (
    "../../SingleElevator/elevatorLogic"
    "../../dataTypes"      
    "../../config"
    "../../SingleElevator/orders"
    "fmt"
    "os"
  )

  const _numElevators int = config.NumElevators  
  const _numFloors int    = config.NumFloors
  const _numOrderButtons int = config.NumOrderButtons

/*
func requests_ClearAtCurrentFloor(Elevator e, onFloor func(Button, int)) Elevator {
  for btn := 0; btn < 4; btn++{
    if (e.Requests[e.Floor][btn]){
      e.Requests[e.Floor][btn] = 0
      if (onFloor){
        onFloor(btn, Floor)
      }
    }
  }
  return e
}*/



func TimeToIdle(elevator dataTypes.ElevatorInfo) int{





  duration := 0

  if (orders.Empty(elevator)){
    return duration
  }

  switch (elevator.State) {
  case dataTypes.S_Disconnected: 
    return 100000000000000
  case dataTypes.S_Idle:
    elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
    if (elevator.CurrentDirection == dataTypes.S_Idle) { //S_Stop
      return duration
    }
  case dataTypes.S_Moving:
    duration += config.ElevatorTravelTimeMs/2
	elevator.Floor += int(elevator.CurrentDirection)


	
  case dataTypes.S_DoorOpen:
	duration += config.DoorOpenTimeMs/2
	//duration += 3000
  }

  // hva sendes inn?




  
  
  a := 0
  for {

	//if (elevator.Floor == 0 && elevator.CurrentDirection == -1 && elevator.State == dataTypes.S_Moving){
		
	//}
	

	if elevatorLogic.ShouldStopHere(elevator){
		elevator.LocalOrders = orders.Execute(elevator)
		duration += config.DoorOpenTimeMs
	}

	if orders.Empty(elevator){
		return duration
	}

	a++

    if a == 98{
      fmt.Println("Should staop on this")
      dataTypes.ElevatorInfoPrint(elevator)
    }

	elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
	if elevator.CurrentDirection != dataTypes.D_Stop{
		elevator.State = dataTypes.S_Moving
	}

	elevator.Floor += int(elevator.CurrentDirection)
	duration += config.ElevatorTravelTimeMs

	

	

    if (a > 100){
      println("Error evig while")
      os.Exit(0)
    }
    

  }
}


/*

    if(elevatorLogic.ShouldStopHere(elevator)){
      fmt.Println("yes")
      dataTypes.ElevatorInfoPrint(elevator)

      elevator := clearAtCurrentFloor(elevator)
      duration += config.DoorOpenTimeMs
      elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
      if elevator.CurrentDirection != dataTypes.D_Stop{
        elevator.State = dataTypes.S_Moving
      }
      if elevator.Floor == 0 || elevator.Floor == _numFloors-1 || orders.Empty(elevator){
        elevator.CurrentDirection = dataTypes.D_Stop
      }
      if(elevator.CurrentDirection == dataTypes.D_Stop){ //S_Stop
        fmt.Println("Yes")
        return duration
      }
    }

    if (a > 100){
      println("Error evig while")
      os.Exit(0)
    }
    

    elevator.Floor += elevatorLogic.FindNextDirection(elevator) //Dir?
	duration += config.ElevatorTravelTimeMs
	
	*/