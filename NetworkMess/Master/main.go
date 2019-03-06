package main

import (
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
	"../dataTypes"
)

/*
type ShortMessage struct {
	floor int
	dir   MotorDirection
	local_orders [3][4]int
	state ElevatorState
}
*/




// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.


func updateTo3(ordersArr* [3][4]int){
	for i := 0; i < 3; i++{
		for j:= 0; j < 4; j++{
			if ordersArr[i][j] == 1{
				ordersArr[i][j] = 3
			}
		}
	}
}


func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	Elevator1 := dataTypes.ElevatorInfo{}
	


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
	helloTx := make(chan dataTypes.LongMessage)
	helloRx := make(chan dataTypes.ShortMessage)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16560, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		for {
			HelloMsg := dataTypes.LongMessage{Elevator1, Elevator1, Elevator1}
			fmt.Println("Sending to elevator:",HelloMsg.Elevator1.Local_orders[1][1])
			helloTx <- HelloMsg
			time.Sleep(1 * time.Second)
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

		case a := <-helloRx:
			//fmt.Printf("Received: %#v\n", a)
			fmt.Println("Recieved from elevator")
			//dataTypes.ElevatorInfoPrint(a.Elevator)
			Elevator1 = a.Elevator
			updateTo3(&Elevator1.Local_orders)
			dataTypes.ElevatorInfoPrint(Elevator1)
			fmt.Println("Should be sending to elevator:",Elevator1.Local_orders[1][1])
		}
	}
}