package main

import (
	"../network/bcast"
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
	 "./costFunction"
)

const _numElevators int = config.NumElevators  
const _numFloors int    = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons

func main() {
	
	priority := os.Args[1:][0] 
    ID, err := strconv.Atoi(priority)
	if err != nil || ID < 1{
		fmt.Println("Not valid  number")
		os.Exit(1)
	}
	
	backupWatchdog := timer.Timer_s{}
	AllElevators := [_numElevators] dataTypes.ElevatorInfo{}

	runAsBackUp := true

	timer.Start(int64(ID*config.WatchdogTimeMs),&backupWatchdog)

	TXElevators       := make(chan [_numElevators] dataTypes.ElevatorInfo)
	RXElevators       := make(chan dataTypes.ElevatorInfo)
	disconnected      := make(chan int)
	TXBackup          := make(chan dataTypes.BackupMessage)
	
	
	watchdogTimedOut  := make(chan bool)
	RXBackup          := make(chan dataTypes.BackupMessage)




	for e := 0; e < _numElevators; e++{
		go elevatorCom.Receive(RXElevators,e, disconnected)
	}
	go bcast.Transmitter(config.BackupRXPort, TXBackup)
	go bcast.Transmitter(config.MasterTXPort, TXElevators)
	go bcast.Receiver(config.BackupRXPort, RXBackup)
	go timer.PollTimer(watchdogTimedOut,&backupWatchdog)



	go func(){
		for{
			time.Sleep(250*time.Millisecond)
			
			if runAsBackUp{
				fmt.Println(" -- Running as BACKUP --")	
			}else{
				fmt.Println(" -- Running as MASTER --")	
			}
			dataTypes.AllElevatorsPrint(AllElevators)
		}
	}()


	/* ----------- Main loop ------------- */

	/* Main loopen er forenkla nå. Om den fungerer dårlig gå tilbake til den som er kommentert bort under, 
		burde egentlig være helt lik. Om den også fungerer dårlig gå tilbake til forrige commit i loggen som heter "må teste for pakketap"
	*/

	for{
		select{
		case MessageFromMaster := <- RXBackup:
			
			if MessageFromMaster.ID != -1 && MessageFromMaster.ID != ID{ // ANNEN MASTER
					
				runAsBackUp = true

				timer.Reset(&backupWatchdog) // hold deg som backup
				TXBackup <- dataTypes.BackupMessage{-1, AllElevators}
			}
		case <-watchdogTimedOut:
			fmt.Println("Watchdog timed out, become master")
			runAsBackUp = false
		case RXelev := <-RXElevators:
			AllElevators = AssignOrdersToElevators(AllElevators, RXelev)
	
			if !runAsBackUp{
				TXElevators <- AllElevators
				TXBackup <- dataTypes.BackupMessage{ID, AllElevators}
			}

		case c := <-disconnected: // elevator number c is disconnected
			AllElevators[c].State = dataTypes.S_Disconnected
			AllElevators = redistributeOrders(AllElevators, c)
		}
	}
		
	
	/*
	for{
		if runAsBackUp{ // BACKUP
			select { // Nå er det ingen dataoverføring mellom master og backup, som er litt weird 
			case MessageFromMaster := <-RXBackup:

				if MessageFromMaster.ID != -1 && MessageFromMaster.ID != ID{ // ANNEN MASTER SOM IKKE ER MEG
					
					// Dette under nødvendig, virker til å lage bugs, test godt på tirsdag 
					//AllElevators = MessageFromMaster.Elevators  

					timer.WatchdogReset(&backupWatchdog) // hold deg som backup
					TXBackup <- dataTypes.BackupMessage{-1, AllElevators}
				}

			case <-watchdogTimedOut:
				fmt.Println("Watchdog timed out, become master")
				runAsBackUp = false

			
			case RXelev := <-RXElevators:
				AllElevators = AssignOrdersToElevators(AllElevators, RXelev)
			
			case c := <-disconnected: // elevator number c is disconnected
				AllElevators[c].State = dataTypes.S_Disconnected
				AllElevators = redistributeOrders(AllElevators, c)
			}		
			

		}else{ // MASTER
			select{
			case a:= <-RXBackup:
				timer.WatchdogReset(&backupWatchdog)
				if a.ID != -1 && a.ID != ID {  // ANNEN MASTER SOM IKKE ER MEG
					fmt.Println("OTHER MASTER CONNECTED BECOME BACKUP")
					runAsBackUp = true
					timer.WatchdogReset(&backupWatchdog)
				}
			case RXelev := <-RXElevators:
					AllElevators = AssignOrdersToElevators(AllElevators, RXelev)
					TXElevators <- AllElevators
					TXBackup <- dataTypes.BackupMessage{ID, AllElevators}


			case c := <-disconnected: // elevator number c is disconnected
				AllElevators[c].State = dataTypes.S_Disconnected
				AllElevators = redistributeOrders(AllElevators, c)
			}
		}
		
	}*/

}

/* If an elevator is disconnected the hall orders of the elevator are given to a connected elevator */
func redistributeOrders(AllElevators [_numElevators] dataTypes.ElevatorInfo, elevatorNumber int)[_numElevators] dataTypes.ElevatorInfo{
	
	ordersToDistribute := AllElevators[elevatorNumber].LocalOrders
	elevatorToDistributeTo := -1
	
	for e := 0; e < _numElevators; e++{
		if AllElevators[e].State != dataTypes.S_Disconnected{
			elevatorToDistributeTo = e		
		}
	}

	if elevatorToDistributeTo == -1{
		fmt.Println("error")
		return AllElevators
	}

	for b := 0; b < _numOrderButtons-1; b++{
		for f:= 0; f < _numFloors; f++{
			if ordersToDistribute[b][f] == dataTypes.O_Handle{
				for e := 0; e < _numElevators; e++{
					if e == elevatorToDistributeTo{
						AllElevators[e].LocalOrders[b][f] = dataTypes.O_Handle
					}else{
						AllElevators[e].LocalOrders[b][f] = dataTypes.O_LightOn
					}
				}
			}
		}
	}
	return AllElevators
}




func lowestCostElevator(AllElevators [_numElevators] dataTypes.ElevatorInfo, NewOrders dataTypes.ElevatorInfo)int{
	lowestCostElevator := -1
	lowestCost := config.InfCost


	for e := 0; e < _numElevators; e++{
		for b := 0; b < _numOrderButtons; b++{ 
			for f:= 0; f < _numFloors; f++{ 

				if NewOrders.LocalOrders[b][f] == dataTypes.O_Received{
					AllElevators[e].LocalOrders[b][f] = dataTypes.O_Handle
				}
			}
		}
		elevatorCost := costFunction.TimeToIdle(AllElevators[e])
		if elevatorCost < lowestCost{
			lowestCostElevator = e
			lowestCost = elevatorCost
		}
	}



	return lowestCostElevator
}




func AssignOrdersToElevators(AllElevators [_numElevators] dataTypes.ElevatorInfo, newOrders dataTypes.ElevatorInfo) [_numElevators] dataTypes.ElevatorInfo {
	
	

	lowestCostElevator := lowestCostElevator(AllElevators, newOrders)

	Elevnumber := newOrders.Number
	AllElevators[Elevnumber].State = newOrders.State
	AllElevators[Elevnumber].Floor = newOrders.Floor
	AllElevators[Elevnumber].CurrentDirection = newOrders.CurrentDirection

	for b := 0; b < _numOrderButtons; b++{
		for f := 0; f <_numFloors; f++{
			for e := 0; e < _numElevators; e++{
				switch newOrders.LocalOrders[b][f]{
				case dataTypes.O_Received:
					for l := 0; l < _numElevators; l++{
						if l == lowestCostElevator{
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_Handle
						}else{
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_LightOn
						}
					}
				case dataTypes.O_Executed:
					if b == dataTypes.BT_Cab{
						AllElevators[Elevnumber].LocalOrders[b][f] = dataTypes.O_Empty
					}else{
						for l := 0; l < _numElevators; l++{
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_Empty
						}
					}
				case dataTypes.O_Handle:
					if b == dataTypes.BT_Cab{ // for cab
						AllElevators[Elevnumber].LocalOrders[b][f] =dataTypes.O_Handle
				}
				default:
					// Do nothing
				}
			}
		}
	}
	return AllElevators
}


