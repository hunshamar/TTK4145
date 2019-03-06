package main

import (
	"./network/bcast"
	//"./network/localip"
	//"./network/peers"

	"fmt"

	"time"

)




type shortMessage struct {
	Message string
	Orders  [3][4]int
}

type longMessage struct{
	Message string
	Orders1  [3][4]int
	Orders2  [3][4]int
	Orders3  [3][4]int
}


var Orders1 = [3][4]int{
	{0, 3, 0, 0}, //up
	{0, 0, 0, 0}, // down
	{0, 0, 0, 0}}  //cab
	
var Orders2 = [3][4]int{
	{0, 0, 0, 0}, //up
	{0, 0, 0, 0}, // down
	{0, 0, 0, 0}}  //cab
			
var Orders3 = [3][4]int{
	{0, 0, 0, 0}, //up
	{0, 0, 0, 0}, // down
	{0, 0, 0, 0}}  //cab
		




func sendOrders(Tx chan<- longMessage, ToTransmit <-chan longMessage){
	sendMEssage :=  <-ToTransmit
	for {
		Tx <- sendMEssage
		time.Sleep(1 * time.Second)
		PrintOrders(sendMEssage.Orders1)
		fmt.Println("Sending")
	}
}


func PrintOrders(arr [3][4]int){
	fmt.Printf("orders_up    : %v \n", arr[0]) 
	fmt.Printf("orders_down  : %v \n", arr[1]) 
	fmt.Printf("orders_panel : %v \n\n", arr[2]) 
}

func main() {

	OrdersStruct:= longMessage{"Hello",Orders1,Orders2,Orders3}


	Rx := make(chan shortMessage)	
	Tx := make(chan longMessage)

	ToTransmit := make(chan longMessage)

	/* ----------------- Start receiving -------------------*/

	
	
	go bcast.Receiver(16569, Rx)
	go bcast.Transmitter(16569, Tx)
	go sendOrders(Tx, ToTransmit)

	
	

	fmt.Printf("\n\n")
	fmt.Println("+-------------------------------------+")
	fmt.Println("| IS BACKUP:  st receiving d1ata |")
	fmt.Println("+-------------------------------------+")
	fmt.Printf("\n")

	runAsBackUp := true
	for runAsBackUp{
		select {
		case a := <-Rx:
			fmt.Println("Recieved data")
			OrdersStruct.Orders1 = a.Orders
			PrintOrders(a.Orders)
		}
		OrdersStruct.Orders1[0][1]++
		ToTransmit <- OrdersStruct
	}
}