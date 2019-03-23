package main

import (
	"../network/bcast"
	"../network/localip"
	"flag"
	"fmt"
	"../timer"
	"time"
	//"log"
	//"os/exec"
	 "os"
	 "strconv"
	 "../dataTypes"
	 "./elevatorCom"
	 "../config"
)

const _numElevators int = config.NumElevators  
const _numFloors int    = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons

func main() {

	priority := os.Args[1:][0] 
    pri, err := strconv.Atoi(priority)
	if err != nil || pri < 1{
		fmt.Println("Not valid number")
		os.Exit(1)
	}

	fmt.Println("seconds:", pri)

	
	AllElevators := dataTypes.AllElevatorInfo{}

	
	TXElevators := make(chan dataTypes.AllElevatorInfo)
	RXElevators := make(chan dataTypes.ElevatorInfo)
	disconnected := make(chan int)

	for e := 1; e <= _numElevators; e++{
	
		go elevatorCom.Receive(RXElevators,e, disconnected)
	}

	TXBackup := make(chan dataTypes.BackupMessage)

	go bcast.Transmitter(config.BackupPort, TXBackup)
	go bcast.Transmitter(config.MasterTXPort, TXElevators)
	//counter := -1
	watchdogTimer1 := timer.Timer_s{}
	timer.WatchdogInit(int64(pri*1000),&watchdogTimer1)

	watchdogTimedOut := make(chan bool)
	RXBackup := make(chan dataTypes.BackupMessage)

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


	runAsBackUp := true


	/* ----------------- Start receiving -------------------*/

	go bcast.Receiver(config.BackupPort, RXBackup)
	go timer.WatchdogPoll(watchdogTimedOut,&watchdogTimer1)




	dataTypes.AllElevatorsPrint(AllElevators)
	fmt.Println("nå skal den være 1, ja")

	go func(){
		time.Sleep(3*time.Second)
		for{
			time.Sleep(1*time.Second)
			dataTypes.AllElevatorsPrint(AllElevators)
			
			if runAsBackUp{
				fmt.Println("IS BACKUP")
			}
		}
	}()

	
	fmt.Println("--- IS BACKUP IS BACKUP ---")
	for{
		if runAsBackUp{ // BACKUP
			select { /* Nå er det ingen dataoverføring mellom master og backup, som er litt weird */
			case a := <-RXBackup:

				if a.ID != -1 && a.ID != pri{ // ANNEN MASTER SOM IKKE ER MEG
					AllElevators = a.Elevators


					timer.WatchdogReset(&watchdogTimer1) // hold deg som backup
					TXBackup <- dataTypes.BackupMessage{-1, AllElevators}
				}
				
	
				
			case <-watchdogTimedOut:
				fmt.Println("Watchdog timed out\n")
				runAsBackUp = false
				dataTypes.AllElevatorsPrint(AllElevators)

			
			case RXelev := <-RXElevators:
				
				if (RXelev.Number != 0){
				

				AllElevators = costFunction(AllElevators, RXelev)
				//AllElevators = clearInvalidOrders(AllElevators)
				}
			}		
			

		}else{ // MASTER
			select{
			case a:= <-RXBackup:
				timer.WatchdogReset(&watchdogTimer1)
				if a.ID != -1 && a.ID != pri {
					fmt.Println("Annen master er koblet til, gtfo")
					fmt.Println("--- IS BACKUP IS BACKUP ---")
					runAsBackUp = true
					timer.WatchdogReset(&watchdogTimer1)
				}
			case RXelev := <-RXElevators:

				if (RXelev.Number != 0){ // Nødvendig? 
				

				AllElevators = costFunction(AllElevators, RXelev)
				//AllElevators = clearInvalidOrders(AllElevators)
				TXElevators <- AllElevators
				TXBackup <- dataTypes.BackupMessage{pri, AllElevators}
				}


			case c := <-disconnected: // elevator number c is disconnected
				AllElevators.Elevators[c-1].State = dataTypes.S_Disconnected
				/* Flytte ordre over til annen heis */ 			

				AllElevators = redistributeOrders(AllElevators, c)
				dataTypes.AllElevatorsPrint(AllElevators)
				//AllElevators = clearInvalidOrders(AllElevators)
			}
		}
		
	}

}

/* If an elevator is disconnected the hall orders of the elevator are given to a connected elevator */
func redistributeOrders(AllElevators dataTypes.AllElevatorInfo, elevatorNumber int)dataTypes.AllElevatorInfo{
	
	ordersToDistribute := AllElevators.Elevators[elevatorNumber-1].LocalOrders
	elevatorToDistributeTo := -1
	
	for e := 0; e < _numElevators; e++{
		if AllElevators.Elevators[e].State != dataTypes.S_Disconnected{
			elevatorToDistributeTo = e		
		}
	}

	if elevatorToDistributeTo == -1{
		fmt.Println("ERROR")
		return AllElevators
	}

	for b := 0; b < _numOrderButtons-1; b++{
		for f:= 0; f < _numFloors; f++{
			if ordersToDistribute[b][f] == 3{
				for e := 0; e < _numElevators; e++{
					if e == elevatorToDistributeTo{
						AllElevators.Elevators[e].LocalOrders[b][f] = 3
					}else{
						AllElevators.Elevators[e].LocalOrders[b][f] = 2
					}
				}
			}
		}
	}
	return AllElevators
}


/* SHOULD NOT BE NECESSARY */ 
func clearInvalidOrders(AllElevators dataTypes.AllElevatorInfo)dataTypes.AllElevatorInfo{
	for b := 0; b < _numOrderButtons; b++{
		for f := 0; f <_numFloors; f++{
			for e := 0; e < _numElevators; e++{
				switch AllElevators.Elevators[e].LocalOrders[b][f]{
				case 2:
					isValid := false
					for l := 0; l < _numElevators; l++{
						if AllElevators.Elevators[l].LocalOrders[b][f] == 3{
							isValid = true
						}		
					}			
					if !isValid{ // If an elevator keeps the light on, but there are no 3 orders on other elevators
						fmt.Println("Is invalid, elevator",e,"Button",b,"Floor",f)
						AllElevators.Elevators[e].LocalOrders[b][f] = 0
					}
				default:
					// Do nothing
				
				}
			}
		}
	}
	return AllElevators
}




func costFunction(AllElevators dataTypes.AllElevatorInfo, newOrders dataTypes.ElevatorInfo) dataTypes.AllElevatorInfo {
	
	Elevnumber := newOrders.Number
	AllElevators.Elevators[Elevnumber -1].State = newOrders.State
	AllElevators.Elevators[Elevnumber -1].Floor = newOrders.Floor
	AllElevators.Elevators[Elevnumber -1].CurrentDirection = newOrders.CurrentDirection

	for b := 0; b < _numOrderButtons; b++{
		for f := 0; f <_numFloors; f++{
			for e := 0; e < _numElevators; e++{
				switch newOrders.LocalOrders[b][f]{
				case 1:
					for l := 0; l < _numElevators; l++{
						if l == Elevnumber-1{
							AllElevators.Elevators[l].LocalOrders[b][f] = 3
						}else{
							AllElevators.Elevators[l].LocalOrders[b][f] = 2
						}
					}
				case -1:
					for l := 0; l < _numElevators; l++{
					
						AllElevators.Elevators[l].LocalOrders[b][f] = 0
						
					}
				case 3:
					AllElevators.Elevators[Elevnumber-1].LocalOrders[b][f] = 3 // For cab calls
				default:
					// Do nothing
				
				}
			}
		}
	}
	return AllElevators
}


