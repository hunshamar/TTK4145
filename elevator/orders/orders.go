
package orders

import "fmt"



var Local_orders =[3][4]int{
	{0,0,0,0},
	{0,0,0,0},
	{0,0,0,0}} 


func PrintOrders(){
	fmt.Printf("orders_up: %v \n", Local_orders[0]) 
	fmt.Printf("orders_down: %v \n", Local_orders[1]) 
	fmt.Printf("orders_panel: %v \n\n", Local_orders[2]) 
}

func Add(buttonType int, floor int){
	Local_orders[buttonType][floor] = 1
}

func Remove(buttonType int, floor int){
	Local_orders[buttonType][floor] = 0
}

