
/*
package main

import "./transmit"

func main(){

	transmit.Transmit()	
}

*/
package main

import (
	"./network/bcast"
	"./network/localip"
	//"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"

)

type shortMessage struct {
	Message string
	Orders    [3][4]int
}

var Local_orders =[3][4]int{
	{0, 0, 0, 0}, //up
	{0, 0, 0, 0}, // down
	{0, 0, 0, 0}}  //cab


func PrintOrders(){
	fmt.Printf("orders_up    : %v \n", Local_orders[0]) 
	fmt.Printf("orders_down  : %v \n", Local_orders[1]) 
	fmt.Printf("orders_panel : %v \n\n", Local_orders[2]) 
}

func sendOrders(Tx chan<- shortMessage, ToTransmit <-chan [3][4]int){
	shortMessage := shortMessage{"From " + "1", <-ToTransmit}
	for {
		Tx <- shortMessage
		time.Sleep(1 * time.Second)
		PrintOrders()
		fmt.Println("Sending")
	}
}

func main() {

	

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()


	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	fmt.Println("Id:",id)

	//peerTxEnable := make(chan bool)
	//go peers.Transmitter(15647, id, peerTxEnable)
	Tx := make(chan shortMessage)
	LOrd := make(chan [3][4]int)

	go bcast.Transmitter(16569, Tx)
	go sendOrders(Tx, LOrd)

	Local_orders[0][1] = 1; 
	

	fmt.Printf("\n\n")
	fmt.Println("+-------------------------------------+")
	fmt.Println("| IS PRIMARY: start transmitting data |")
	fmt.Println("+-------------------------------------+")
	fmt.Printf("\n")
	for {
		LOrd <- Local_orders
	}
}
