
package buttons

import "../elevio"
import "../../dataTypes"
import "fmt"
import "time"

const _pollRate = 20 * time.Millisecond

func MirrorOrders(elevatorInfo <-chan dataTypes.ElevatorInfo ){
	for{

		elevator := <- elevatorInfo

		for floor := 0; floor <= 3; floor++{
			if elevator.LocalOrders[0][floor] >= 2 {
				elevio.SetButtonLamp(dataTypes.BT_HallUp, floor, true)
			} else{
			elevio.SetButtonLamp(dataTypes.BT_HallUp, floor, false)
			}
		}
		for floor := 0; floor <= 3; floor++{
			if elevator.LocalOrders[1][floor] >= 2 {
				elevio.SetButtonLamp(dataTypes.BT_HallDown, floor, true)
			} else{
			elevio.SetButtonLamp(dataTypes.BT_HallDown, floor, false)
			}
		}
		for floor := 0; floor <= 3; floor++{
			if elevator.LocalOrders[2][floor] >= 2 {
				elevio.SetButtonLamp(dataTypes.BT_Cab, floor, true)
			} else{
			elevio.SetButtonLamp(dataTypes.BT_Cab, floor, false)
			}
		}
		time.Sleep(_pollRate)
	}
}



func Pr(){
	fmt.Println("a")
}


