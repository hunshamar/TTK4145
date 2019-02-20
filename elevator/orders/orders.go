
package orders

import "fmt"
import ""

var O_hall_up = [4]int{0,0,0,0}
var O_hall_down = [4]int{0,0,0,0}
var O_cab = [4]int{0,0,0,0}

var orders =[3][4]{
	{0,0,0,0},
	{0,0,0,0},
	{0,0,0,0}
} 

orders[button][floor]


func PrintOrders(){
	fmt.Printf("orders_up: %v \n", O_hall_up) 
	fmt.Printf("orders_down: %v \n", O_hall_down) 
	fmt.Printf("orders_panel: %v \n\n", O_cab) 
}

func Add(buttonType int, floor int){
	if buttonType == 0{
		O_hall_up[floor] = 1
	}
	if buttonType == 1{
		O_hall_down[floor] = 1
	}
	if buttonType == 2{
		O_cab[floor] = 1
	}
}

