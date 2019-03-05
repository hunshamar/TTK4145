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

type BackupMsg struct {
	Message string
	Iter    int
}

func spawnBackup(){
	cmd := exec.Command("gnome-terminal", "--window", "-x","go", "run", "primary.go") 
	cmd.Stderr = os.Stderr


	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}
}

func main() {

	counter := 0
	watchdogTimer1 := timer.Timer_s{}
	timer.WatchdogInit(3000,&watchdogTimer1)

	watchdogTimedOut := make(chan bool)
	helloRx := make(chan BackupMsg)

	/* ----------------- Start receiving -------------------*/

	go bcast.Receiver(16569, helloRx)
	go timer.WatchdogPoll(watchdogTimedOut,&watchdogTimer1)

	fmt.Printf("\n\n")
	fmt.Println("+-------------------------------------+")
	fmt.Println("| IS BACKUP:  starting receiving data |")
	fmt.Println("+-------------------------------------+")
	fmt.Printf("\n")

	runAsBackUp := true
	for runAsBackUp{
		select {
		case a := <-helloRx:
			counter = a.Iter
			timer.WatchdogReset(&watchdogTimer1)
			
		case <-watchdogTimedOut:
			fmt.Printf("Watchdog timed out\n")
			runAsBackUp = false
		}
	}
	fmt.Println("PRIMARY has crashed, become PRIMARY and open new BACKUP")



	
	/* ----------------- Spawn backup -------------------*/
	spawnBackup()

	
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
	helloTx := make(chan BackupMsg)
	
	go bcast.Transmitter(16569, helloTx)

	go func() {
		BackupMsg := BackupMsg{"From " + id, counter}
		for {
			BackupMsg.Iter++
			helloTx <- BackupMsg
			time.Sleep(1 * time.Second)
			fmt.Println("Count:", BackupMsg.Iter)
		}
	}()
	fmt.Printf("\n\n")
	fmt.Println("+-------------------------------------+")
	fmt.Println("| IS PRIMARY: start transmitting data |")
	fmt.Println("+-------------------------------------+")
	fmt.Printf("\n")
	for {
		// do nothing, looping as primary
	}

}

