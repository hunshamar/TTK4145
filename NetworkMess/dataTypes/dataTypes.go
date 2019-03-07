
package dataTypes

import "fmt"

type ElevatorState int 
const (
	Idle = 0
	Moving = 1
	DoorOpen = 2
)

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
	Floor int
	CurrentDirection   Direction
	LocalOrders [3][4]int
	State ElevatorState
}

type ShortMessage struct {
	Elevator ElevatorInfo
}

type LongMessage struct{
	Elevator1 ElevatorInfo
	Elevator2 ElevatorInfo
	Elevator3 ElevatorInfo
}

type HelloMsg struct {
	Message string
	Iter    int
}

func ElevatorInfoPrint(EI ElevatorInfo){
	fmt.Println("----------------------")
	fmt.Println("Floor:",EI.Floor) 
	fmt.Println("Direction:",EI.CurrentDirection) 
	fmt.Println("Orders:",EI.LocalOrders) 
	fmt.Println("State:",EI.State) 
	fmt.Println("----------------------")
}