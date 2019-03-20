
package orders

import "fmt"
import "../../config"
import "../../dataTypes"

func PrintOrders(elevator dataTypes.ElevatorInfo){
	fmt.Printf("orders_up    : %v \n", elevator.LocalOrders[0]) 
	fmt.Printf("orders_down  : %v \n", elevator.LocalOrders[1]) 
	fmt.Printf("orders_panel : %v \n\n", elevator.LocalOrders[2]) 
}

const _numFloors int= config.NumFloors
const _numOrderButtons int = config.NumOrderButtons

func Add(elevator dataTypes.ElevatorInfo, button dataTypes.ButtonType, floor int ) [_numOrderButtons][_numFloors]int{
	switch button{
	case dataTypes.BT_Cab:
		elevator.LocalOrders[button][floor] = dataTypes.O_Handle
	case dataTypes.BT_HallUp:
		fallthrough
	case dataTypes.BT_HallDown:
		elevator.LocalOrders[button][floor] = dataTypes.O_Received 
	}
	return elevator.LocalOrders
}

func Execute(elevator dataTypes.ElevatorInfo) [_numOrderButtons][_numFloors]int{
	
	localOrders := elevator.LocalOrders
	for buttonType := 0; buttonType < _numOrderButtons; buttonType++{ 
		if (localOrders[buttonType][elevator.Floor] == dataTypes.O_Handle){ 
			localOrders[buttonType][elevator.Floor] = dataTypes.O_Executed
		}
	}
	return localOrders
}


func Above(elevator dataTypes.ElevatorInfo,floor int)bool{
	if (floor == _numFloors -1){
		return false
	}
	for b := 0; b < _numOrderButtons; b++{
		for f := floor+1; f < _numFloors; f++{
			if elevator.LocalOrders[b][f] == dataTypes.O_Handle{
				return true
			}
		}
	}
	return false
}


func Below(elevator dataTypes.ElevatorInfo,floor int)bool{
	if (floor == 0){
		return false
	}

	for b := 0; b < 3; b++{
		for f := 0; f < floor; f++{
			if elevator.LocalOrders[b][f] == dataTypes.O_Handle{
				return true
			}
		}
	}

	return false
}


func Empty(elevator dataTypes.ElevatorInfo)bool{
	return !Above(elevator,0) && !Below(elevator,_numFloors-1)
}


func OnFloor(elevator dataTypes.ElevatorInfo,floor int)bool{
	return Above(elevator,floor-1) && Below(elevator,floor+1)
}


