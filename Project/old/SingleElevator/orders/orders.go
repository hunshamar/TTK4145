
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


func Execute(elevator dataTypes.ElevatorInfo) [_numOrderButtons][_numFloors]int{
	
	localOrders := elevator.LocalOrders
	for buttonType := 0; buttonType < 2; buttonType++{ // hall buttons
		if (localOrders[buttonType][elevator.Floor] == dataTypes.O_Handle){ 
			localOrders[buttonType][elevator.Floor] = dataTypes.O_Executed
		}
	}
	if (localOrders[dataTypes.BT_Cab][elevator.Floor] == dataTypes.O_Handle){
		localOrders[dataTypes.BT_Cab][elevator.Floor] = dataTypes.O_Executed // husk å gjøre lettere om funker
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
				fmt.Println("Jepp")
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


