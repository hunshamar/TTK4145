
package orders

import "fmt"
import "../elevio"

var Local_orders =[3][4]bool{
	{false, false, false, false}, //up
	{false, false, false, false}, // down
	{false, false, false, false}}  //cab


func PrintOrders(){
	fmt.Printf("orders_up: %v \n", Local_orders[0]) 
	fmt.Printf("orders_down: %v \n", Local_orders[1]) 
	fmt.Printf("orders_panel: %v \n\n", Local_orders[2]) 
}

func Add(buttonType int, floor int){
	Local_orders[buttonType][floor] = true
}

func Remove(floor int){
	for buttonType := 0; buttonType <= 2; buttonType++{
		Local_orders[buttonType][floor] = false
	}
}

func ExecutableOnFloor(s int, floor int) bool{
	state := s
	switch state{
	case 1:
		if (Local_orders[elevio.BT_Cab][floor] || Local_orders[elevio.BT_HallUp][floor] || (Local_orders[elevio.BT_HallDown][floor] && !Above(floor) )){
			return true
		}
	case -1:
		if (Local_orders[elevio.BT_Cab][floor] || Local_orders[elevio.BT_HallDown][floor] || (Local_orders[elevio.BT_HallUp][floor] && !Below(floor) )){
			return true
		}
	case 0:
		if (Local_orders[elevio.BT_Cab][floor] || Local_orders[elevio.BT_HallDown][floor] || Local_orders[elevio.BT_HallUp][floor] ){
			return true
		}
	case 2:
		if (Local_orders[elevio.BT_Cab][floor] || Local_orders[elevio.BT_HallDown][floor] || Local_orders[elevio.BT_HallUp][floor] ){
			return true
		}
	}
	return false
}

func Above(floor int)bool{
	if (floor == 3){
		fmt.Println("ERROR")
		return false
	}


	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor+1; floor <= 3; floor++{
			if Local_orders[buttonType][floor]{
				return true
			}
		}
	}

	return false
}


func Below(floor int)bool{
	if (floor == 0){
		fmt.Println("ERROR")
		return false
	}

	for buttonType := 0; buttonType <= 2; buttonType++{
		for floor := floor-1; floor >= 0; floor--{
			if Local_orders[buttonType][floor]{
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


