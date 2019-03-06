
package orders

import "fmt"
import "../elevio"

var Local_orders =[3][4]int{
	{0, 0, 0, 0}, //up
	{0, 0, 0, 0}, // down
	{0, 0, 0, 0}}  //cab


func PrintOrders(){
	fmt.Printf("orders_up    : %v \n", Local_orders[0]) 
	fmt.Printf("orders_down  : %v \n", Local_orders[1]) 
	fmt.Printf("orders_panel : %v \n\n", Local_orders[2]) 
}

func Add(buttonType int, floor int, value int){
	Local_orders[buttonType][floor] = value
}

func Remove(floor int){
	for buttonType := 0; buttonType <= 2; buttonType++{
		Local_orders[buttonType][floor] = 0
	}
}

func ExecutableOnFloor(s int, floor int) bool{
	state := s
	switch state{
	case 1:
		if (Local_orders[elevio.BT_Cab][floor] ==  3 || Local_orders[elevio.BT_HallUp][floor] ==  3 || (Local_orders[elevio.BT_HallDown][floor] ==  3 && !Above(floor) )){
			return true
		}
	case -1:
		if (Local_orders[elevio.BT_Cab][floor] ==  3 || Local_orders[elevio.BT_HallDown][floor] ==  3 || (Local_orders[elevio.BT_HallUp][floor] ==  3 && !Below(floor) )){
			return true
		}
	case 0:
		if (Local_orders[elevio.BT_Cab][floor] ==  3 || Local_orders[elevio.BT_HallDown][floor] ==  3 || Local_orders[elevio.BT_HallUp][floor] ==  3 ){
			return true
		}
	case 2:
		if (Local_orders[elevio.BT_Cab][floor] ==  3 || Local_orders[elevio.BT_HallDown][floor] ==  3 || Local_orders[elevio.BT_HallUp][floor] ==  3 ){
			return true
		}
	}
	return false
}

func Above(floor int)bool{
	if (floor == 3){
		return false
	}


	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor+1; floor <= 3; floor++{
			if Local_orders[buttonType][floor] == 3{
				return true
			}
		}
	}

	return false
}


func Below(floor int)bool{
	if (floor == 0){
		return false
	}

	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor-1; floor >= 0; floor--{
			if Local_orders[buttonType][floor] == 3{
				return true
			}
		}
	}

	return false
}


func Empty()bool{
	return !Above(0) && !Below(3)
}

func OnFloor(floor int)bool{
	return Above(floor-1) && Below(floor+1)
}


