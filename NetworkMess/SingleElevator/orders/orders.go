
package orders

import "fmt"
import "../../dataTypes"

func PrintOrders(elevator dataTypes.ElevatorInfo){
	fmt.Printf("orders_up    : %v \n", elevator.LocalOrders[0]) 
	fmt.Printf("orders_down  : %v \n", elevator.LocalOrders[1]) 
	fmt.Printf("orders_panel : %v \n\n", elevator.LocalOrders[2]) 
}


func ExecutableOnFloor(elevator dataTypes.ElevatorInfo,s int, floor int) bool{
	state := s
	switch state{
	case 1:
		if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  3 || (elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  3 && !Above(elevator,floor) )){
			return true
		}
	case -1:
		if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  3 || (elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  3 && !Below(elevator,floor) )){
			return true
		}
	case 0:
		if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  3 ){
			return true
		}
	case 2:
		if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  3 || elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  3 ){
			return true
		}
	}
	return false
}


func Above(elevator dataTypes.ElevatorInfo,floor int)bool{
	if (floor == 3){
		return false
	}
	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor+1; floor <= 3; floor++{
			if elevator.LocalOrders[buttonType][floor] == 3{
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

	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor-1; floor >= 0; floor--{
			if elevator.LocalOrders[buttonType][floor] == 3{
				return true
			}
		}
	}

	return false
}


func Empty(elevator dataTypes.ElevatorInfo)bool{
	return !Above(elevator,0) && !Below(elevator,3)
}

func OnFloor(elevator dataTypes.ElevatorInfo,floor int)bool{
	return Above(elevator,floor-1) && Below(elevator,floor+1)
}


