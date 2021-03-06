package main

import (
	"fmt"
	"time"
	"os"
	"../network/bcast"
	"../timer"
	"strconv"
	"../config"
	"../dataTypes"
	"./costFunction"
	"./elevatorCom"
)

const _numElevators int    = config.NumElevators
const _numFloors int       = config.NumFloors
const _numOrderButtons int = config.NumOrderButtons



func main() {

	priority := os.Args[1:][0]
	ID, err := strconv.Atoi(priority)
	if err != nil || ID < 1 {
		fmt.Println("Not valid  number")
		os.Exit(1)
	}

	backupWatchdog := timer.Timer_s{}
	AllElevators := [_numElevators]dataTypes.ElevatorInfo{}
	state := dataTypes.MS_Master

	backupWatchdog.Start(ID * config.WatchdogTimeMs)

	TXElevators := make(chan [_numElevators]dataTypes.ElevatorInfo)
	RXElevators := make(chan dataTypes.ElevatorInfo)
	disconnected := make(chan int)
	TXBackup := make(chan dataTypes.BackupMessage)
	backupWatchdogTimedOut := make(chan bool)
	RXBackup := make(chan dataTypes.BackupMessage)

	for e := 0; e < _numElevators; e++ {
		go elevatorCom.ReceiveAndCheckStatus(RXElevators, e, disconnected)
	}
	go bcast.Transmitter(config.BackupRXPort, TXBackup)
	go bcast.Transmitter(config.MasterTXPort, TXElevators)
	go bcast.Receiver(config.BackupRXPort, RXBackup)
	go backupWatchdog.PollTimer(backupWatchdogTimedOut)

	go func() { // Print main status 4 hz
		for {
			time.Sleep(250 * time.Millisecond)
			if state == dataTypes.MS_Backup {
				fmt.Println(" -- Running as BACKUP --")
			} else {
				fmt.Println(" -- Running as MASTER --")
			}
			dataTypes.AllElevatorsPrint(AllElevators)
		}
	}()

	/* ----------- Main loop ------------- */

	for {
		select {
		case MessageFromMaster := <-RXBackup:

			if MessageFromMaster.ID != -1 && MessageFromMaster.ID != ID { // OTHER MASTER

				state = dataTypes.MS_Backup
				AllElevators = MessageFromMaster.Elevators

				backupWatchdog.Reset() // Stay as backup
				TXBackup <- dataTypes.BackupMessage{-1, AllElevators}
			}
		case <-backupWatchdogTimedOut:
			fmt.Println("Watchdog timed out, become master")
			state = dataTypes.MS_Master
		case RXelev := <-RXElevators:

			AllElevators = AssignOrdersToElevators(AllElevators, RXelev)

			if state == dataTypes.MS_Master {
				TXElevators <- AllElevators
				TXBackup <- dataTypes.BackupMessage{ID, AllElevators}
			}

		case c := <-disconnected: // elevator number c is disconnected
			AllElevators[c].State = dataTypes.S_Disconnected
			AllElevators = redistributeOrders(AllElevators, c)
		}
	}
}

/* If an elevator is disconnected the hall orders of the elevator are given to a connected elevator */
func redistributeOrders(AllElevators [_numElevators]dataTypes.ElevatorInfo, elevatorNumber int) [_numElevators]dataTypes.ElevatorInfo {

	ordersToDistribute := AllElevators[elevatorNumber].LocalOrders
	elevatorToDistributeTo := -1

	for e := 0; e < _numElevators; e++ {
		if AllElevators[e].State != dataTypes.S_Disconnected {
			elevatorToDistributeTo = e
		}
	}

	if elevatorToDistributeTo == -1 {
		fmt.Println("error")
		return AllElevators
	}

	for b := 0; b < _numOrderButtons-1; b++ {
		for f := 0; f < _numFloors; f++ {
			if ordersToDistribute[b][f] == dataTypes.O_Handle {
				for e := 0; e < _numElevators; e++ {
					if e == elevatorToDistributeTo {
						AllElevators[e].LocalOrders[b][f] = dataTypes.O_Handle
					} else {
						AllElevators[e].LocalOrders[b][f] = dataTypes.O_LightOn
					}
				}
			}
		}
	}
	return AllElevators
}

func lowestCostElevator(AllElevators [_numElevators]dataTypes.ElevatorInfo, newOrders dataTypes.ElevatorInfo) int {
	lowestCostElevator := newOrders.Number //-1
	lowestCost := config.InfCost

	for e := 0; e < _numElevators; e++ {
		for b := 0; b < _numOrderButtons; b++ {
			for f := 0; f < _numFloors; f++ {

				if newOrders.LocalOrders[b][f] == dataTypes.O_Received {
					AllElevators[e].LocalOrders[b][f] = dataTypes.O_Handle
				}
			}
		}
		elevatorCost := costFunction.TimeToIdle(AllElevators[e])
		if elevatorCost < lowestCost {
			lowestCostElevator = e
			lowestCost = elevatorCost
		}
	}
	return lowestCostElevator
}

func AssignOrdersToElevators(AllElevators [_numElevators]dataTypes.ElevatorInfo, newOrders dataTypes.ElevatorInfo) [_numElevators]dataTypes.ElevatorInfo {

	Elevnumber := newOrders.Number
	AllElevators[Elevnumber].State = newOrders.State
	AllElevators[Elevnumber].Floor = newOrders.Floor
	AllElevators[Elevnumber].CurrentDirection = newOrders.CurrentDirection
	lowestCostElevator := lowestCostElevator(AllElevators, newOrders)

	for b := 0; b < _numOrderButtons; b++ {
		for f := 0; f < _numFloors; f++ {
			for e := 0; e < _numElevators; e++ {
				switch newOrders.LocalOrders[b][f] {
				case dataTypes.O_Received:
					for l := 0; l < _numElevators; l++ {
						if l == lowestCostElevator {
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_Handle
						} else {
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_LightOn
						}
					}
				case dataTypes.O_Executed:
					if b == dataTypes.BT_Cab {
						AllElevators[Elevnumber].LocalOrders[b][f] = dataTypes.O_Empty
					} else {
						for l := 0; l < _numElevators; l++ {
							AllElevators[l].LocalOrders[b][f] = dataTypes.O_Empty
						}
					}
				case dataTypes.O_Handle:
					if b == dataTypes.BT_Cab { // for cab
						AllElevators[Elevnumber].LocalOrders[b][f] = dataTypes.O_Handle
					}
				default:
					// Do nothing
				}
			}
		}
	}
	return AllElevators
}
