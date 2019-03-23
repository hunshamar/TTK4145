package elevatorLogic

import(
	"../../dataTypes"
	"fmt"
	"../orders"
)



func ShouldStopHere(elevator dataTypes.ElevatorInfo) bool{

	floor := elevator.Floor

	if elevator.State == dataTypes.S_Disconnected{
		return false
	}

	switch elevator.State{
	case dataTypes.S_Moving:

		if (elevator.CurrentDirection == dataTypes.D_Up){
			if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  dataTypes.O_Handle || elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  dataTypes.O_Handle || (elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  dataTypes.O_Handle && !orders.Above(elevator,floor) )){
				fmt.Println("Yes stop here ")
				return true
			}
		}
		if (elevator.CurrentDirection == dataTypes.D_Down){
			if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  dataTypes.O_Handle || elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  dataTypes.O_Handle || (elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  dataTypes.O_Handle && !orders.Below(elevator,floor) )){
				return true
			}
		}
		
	case dataTypes.S_Idle:
		fallthrough
	case dataTypes.S_DoorOpen:
		if (elevator.LocalOrders[dataTypes.BT_Cab][floor] ==  dataTypes.O_Handle || elevator.LocalOrders[dataTypes.BT_HallDown][floor] ==  dataTypes.O_Handle || elevator.LocalOrders[dataTypes.BT_HallUp][floor] ==  dataTypes.O_Handle ){
			return true
		}
	}
	return false
} 

func FindNextDirection(elevator dataTypes.ElevatorInfo) dataTypes.Direction{
	if orders.Empty(elevator){
		return dataTypes.D_Stop
	}
	switch elevator.CurrentDirection{
		case dataTypes.D_Stop:
			if orders.Above(elevator, elevator.Floor){
				return dataTypes.D_Up
			}
			if orders.Below(elevator, elevator.Floor){
				return dataTypes.D_Down
			}

		case dataTypes.D_Up:
			if orders.Above(elevator, elevator.Floor){
				return dataTypes.D_Up
			}else{
				return dataTypes.D_Down //Snu, finnes ordre andre veien
			}

		case dataTypes.D_Down:
			if orders.Below(elevator, elevator.Floor){
				return dataTypes.D_Down
			}else{
				return dataTypes.D_Up //Snu, finnes ordre andre veien
			}
	}
	return dataTypes.D_Stop	
}
