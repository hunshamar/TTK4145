
package dataTypes


type ElevatorState int 
const (
	MovingDown = -1
	Idle = 0
	MovingUp = 1	
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
