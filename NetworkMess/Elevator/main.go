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
	"./masterCom"
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


func main() {

	var a =[3][4]int{
		{0, 0, 0, 0}, //up
		{0, 1, 0, 0}, // down
		{0, 0, 0, 0}}  //cab

	Info := dataTypes.ElevatorInfo{
		Floor: 0,
		Dir: dataTypes.MD_Stop,
		Local_orders: a,
		State: dataTypes.Idle}




	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
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

	infoToMaster := make(chan dataTypes.ElevatorInfo)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan dataTypes.ShortMessage)
	helloRx := make(chan dataTypes.LongMessage)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16561, helloTx)
	go bcast.Receiver(16569, helloRx)

	// The example message. We just send one of these every second.


	go masterCom.Transmit(helloTx, infoToMaster)
	/*go func() {
		HelloMsg := dataTypes.ShortMessage{Info}

		for {
			helloTx <- HelloMsg
			time.Sleep(1 * time.Second)
		}
	}()*/

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Println("Recieved from master")

			fmt.Println("\nElevator1:")
			dataTypes.ElevatorInfoPrint(a.Elevator1)
			fmt.Println("\nElevator2:")
			dataTypes.ElevatorInfoPrint(a.Elevator2)
			fmt.Println("\nElevator3:")
			dataTypes.ElevatorInfoPrint(a.Elevator3)


		}
		infoToMaster <- Info
	}
}