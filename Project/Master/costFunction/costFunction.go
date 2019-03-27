package costFunction

import (
  "../../SingleElevator/elevatorLogic"
  "../../dataTypes"      
  "../../config"
  "../../SingleElevator/orders"
)

const _numElevators int = config.NumElevators  
const _numFloors int    = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons

func TimeToIdle(elevator dataTypes.ElevatorInfo) int{

  duration := 0

  if (orders.Empty(elevator)){
    return duration
  }

  switch (elevator.State) {
  case dataTypes.S_Disconnected: 
    return config.InfCost
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
  }
  
  for {

    if elevatorLogic.ShouldStopHere(elevator){
      elevator.LocalOrders = orders.Execute(elevator)
      duration += config.DoorOpenTimeMs
    }

    if orders.Empty(elevator){
      return duration
    }

    elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
    if elevator.CurrentDirection != dataTypes.D_Stop{
      elevator.State = dataTypes.S_Moving
    }

    elevator.Floor += int(elevator.CurrentDirection)
    duration += config.ElevatorTravelTimeMs
  }
}

