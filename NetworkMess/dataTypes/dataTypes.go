
package dataTypes


type ElevatorState int 
const (
	Moving = 0
	Idle = 1
	DoorOpen = 2
)

type MotorDirection int
const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int
const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ShortMessage struct {
	floor int
	dir   MotorDirection
	local_orders [3][4]int
	state ElevatorState
}

type LongMessage struct{
	Elevator1 ShortMessage
	Elevator2 ShortMessage
	Elevator3 ShortMessage
}

type HelloMsg struct {
	Message string
	Iter    int
}