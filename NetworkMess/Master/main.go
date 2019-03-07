package main

import (
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"flag"
	"fmt"
	"os"
	//"time"
	"../dataTypes"
	"./elevatorCom"
	"time"
)

/*
type ShortMessage struct {
	floor int
	dir   MotorDirection
	LocalOrders [3][4]int
	state ElevatorState
}
*/




// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.


func updateTo3(ordersArr1* [3][4]int,ordersArr2* [3][4]int,ordersArr3* [3][4]int, elevator int){

	switch elevator{
	case 1:
		for i := 0; i < 3; i++{
			for j:= 0; j < 4; j++{
				if ordersArr1[i][j] == 1{
					ordersArr1[i][j] = 3
					ordersArr2[i][j] = 2
					ordersArr3[i][j] = 2
				}
				if ordersArr1[i][j] == -1{
					ordersArr1[i][j] = 0
					ordersArr2[i][j] = 0
					ordersArr3[i][j] = 0
				}
			}
		}
	case 2:
		for i := 0; i < 3; i++{
			for j:= 0; j < 4; j++{
				if ordersArr2[i][j] == 1{
					ordersArr1[i][j] = 2
					ordersArr2[i][j] = 3
					ordersArr3[i][j] = 2
				}
				if ordersArr2[i][j] == -1{
					ordersArr1[i][j] = 0
					ordersArr2[i][j] = 0
					ordersArr3[i][j] = 0
				}
			}
		}
	case 3:
		for i := 0; i < 3; i++{
			for j:= 0; j < 4; j++{
				if ordersArr3[i][j] == 1{
					ordersArr1[i][j] = 2
					ordersArr2[i][j] = 2
					ordersArr3[i][j] = 3
				}
				if ordersArr3[i][j] == -1{
					ordersArr1[i][j] = 0
					ordersArr2[i][j] = 0
					ordersArr3[i][j] = 0
				}
			}
		}
	}
}


func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	elevator1 := dataTypes.ElevatorInfo{}
	elevator2 := dataTypes.ElevatorInfo{}
	elevator3 := dataTypes.ElevatorInfo{}



	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	TXToAll := make(chan dataTypes.LongMessage)
	RxElev1 := make(chan dataTypes.ShortMessage)
	RxElev2 := make(chan dataTypes.ShortMessage)
	RxElev3 := make(chan dataTypes.ShortMessage)


	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, TXToAll)
	go bcast.Receiver(16561, RxElev1)
	go bcast.Receiver(16562, RxElev2)
	go bcast.Receiver(16563, RxElev3)

	// The example message. We just send one of these every second.
	infoToElevators := make(chan dataTypes.LongMessage)

	go elevatorCom.Transmit(TXToAll, infoToElevators)

	go func(){
		
		for{
			infoToElevators <- dataTypes.LongMessage{elevator1, elevator2, elevator3}
		time.Sleep(250 * time.Millisecond)
		dataTypes.ElevatorInfoPrint(elevator3)
		}
	 }()


	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-RxElev1:
			elevator1 = a.Elevator
			if elevator1.LocalOrders[1][1] == 1{
				fmt.Println("JAJAJA \n\n\n\n\n\n\n\n\n\n\n\n\n\n")
			}
			updateTo3(&elevator1.LocalOrders,&elevator2.LocalOrders,&elevator3.LocalOrders,1)
			if elevator1.LocalOrders[1][1] == 3{
				fmt.Println("NEINEINEI \n\n\n\n\n\n\n\n\n\n\n\n\n\n")
			}

		case b := <-RxElev2:
			elevator2 = b.Elevator
			if elevator2.LocalOrders[1][1] == 1{
				fmt.Println("JAJAJA \n\n\n\n\n\n\n\n\n\n\n\n\n\n")
			}
			updateTo3(&elevator1.LocalOrders,&elevator2.LocalOrders,&elevator3.LocalOrders,2)
			if elevator2.LocalOrders[1][1] == 3{
				fmt.Println("NEINEINEI \n\n\n\n\n\n\n\n\n\n\n\n\n\n")
			}
		case c := <-RxElev3:
			//fmt.Printf("Received: %#v\n", a)
			fmt.Println("Recieved from elevator")
			//dataTypes.ElevatorInfoPrint(a.Elevator)
			elevator3 = c.Elevator
			updateTo3(&elevator1.LocalOrders,&elevator2.LocalOrders,&elevator3.LocalOrders,3)
			dataTypes.ElevatorInfoPrint(elevator3)
		}
		
	}
}