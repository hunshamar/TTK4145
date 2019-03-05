package main

import (
	"./network/bcast"
	"./network/localip"
	//"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
	"./timer"
	"log"
	"os/exec"
)

type HelloMsg struct {
	Message string
	Iter    int
}

func main() {

	counter := 0

	timer.WatchdogInit(3000)

	//peerUpdateCh := make(chan peers.PeerUpdate)

	//go peers.Receiver(15647, peerUpdateCh)

	watchdog := make(chan bool)
	helloRx := make(chan HelloMsg)

	/* ----------------- Start receiving -------------------*/

	go bcast.Receiver(16569, helloRx)
	go timer.WatchdogPoll(watchdog)

	fmt.Println("Started recieving (backup)")
	
	exit := false
	for !exit{
		select {
		/*case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New) 	
			fmt.Printf("  Lost:     %q\n", p.Lost)*/

		case a := <-helloRx:
			//fmt.Printf("Received: %#v\n", a)
			counter = a.Iter


			timer.WatchdogReset()
		case w := <-watchdog:
			fmt.Printf("Watchdog sees: %#v\n", w)
			exit = true
		}
	}
	fmt.Println("Primary has crashed, open backup")


	/* ----------------- Spawn backup -------------------*/
	cmd := exec.Command("gnome-terminal", "--window", "-x","go", "run", "primary.go") 
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr


	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}

	/* ----------------- Start transmiting -------------------*/
	
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

	//peerTxEnable := make(chan bool)
	//go peers.Transmitter(15647, id, peerTxEnable)
	helloTx := make(chan HelloMsg)
	
	go bcast.Transmitter(16569, helloTx)

	go func() {
		helloMsg := HelloMsg{"Fra main2.go from " + id, counter}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
			fmt.Println("Count:", helloMsg.Iter)
		}
	}()

	fmt.Println("Started sending as primary")
	for {
		
	}

}

