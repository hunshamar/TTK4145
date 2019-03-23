
package dataTypes

import "fmt"
import "../config"


const _numFloors int = config.NumFloors
const _numElevators int = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons

type ElevatorState int 
const (
	S_Idle = 0
	S_Moving = 1
	S_DoorOpen = 2
	S_Disconnected = 3
)

type OrderStatus int
const (
	O_Empty = 0
	O_Executed = -1
	O_Received = 1
	O_LightOn = 2 
	O_Handle = 3
)

func StateToString(s ElevatorState) string{
	switch s{
	case S_Idle:
		return "Idle"
	case S_Moving:
		return "Moving"
	case S_DoorOpen:
		return "Door Open"
	case S_Disconnected:
		return "Disconnected"
	}
	return "Not defined"
}

type Direction int
const (
	D_Up                  = 1
	D_Down                = -1
	D_Stop                = 0
)

type ButtonType int
const (
	BT_HallUp              = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)




type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ElevatorInfo struct{
	HardwareFunctioning bool
	Number int
	Floor int
	CurrentDirection   Direction
	LocalOrders [_numOrderButtons][_numFloors]int
	State ElevatorState
}

type AllElevatorInfo struct{
	Elevators [_numElevators] ElevatorInfo
}



type BackupMessage struct{
	ID int
	Elevators [_numElevators] ElevatorInfo
}

type HelloMsg struct {
	Message string
	Iter    int
}

func OrdersPrint(EI ElevatorInfo){
	fmt.Print(EI.LocalOrders)
	fmt.Print("  N: ", EI.Number)
	fmt.Print("  F: ", EI.Floor)
	fmt.Print("  D: ", EI.CurrentDirection)
	fmt.Print("  S: ", StateToString(EI.State) )

	fmt.Println()
}

func ElevatorInfoPrint(EI ElevatorInfo){
	fmt.Println("----------------------")
	fmt.Println("Floor:",EI.Floor) 
	fmt.Println("Direction:",EI.CurrentDirection) 
	fmt.Println("Orders:",EI.LocalOrders) 
	fmt.Println("State:", StateToString(EI.State)) 
	fmt.Println("----------------------")
}

func AllElevatorsPrint(Elevators [_numElevators] ElevatorInfo){
	fmt.Println("             | H  up | | H  dn | |  Cab   |")

	for e := 0; e < _numElevators; e++{
		fmt.Print("Elevator ", e, "  ")
		OrdersPrint(Elevators[e])
	}
}
