
package buttons

import "../orders"
import "../elevio"
import "fmt"
import "time"

const _pollRate = 20 * time.Millisecond

func MirrorOrders(){
	for{
		
		for floor := 0; floor <= 3; floor++{
			if orders.Local_orders[0][floor] {
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, true)
			} else{
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			}
		}
		for floor := 0; floor <= 3; floor++{
			if orders.Local_orders[1][floor] {
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, true)
			} else{
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			}
		}
		for floor := 0; floor <= 3; floor++{
			if orders.Local_orders[2][floor] {
				elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
			} else{
			elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
			}
		}
		
		
		
		
			time.Sleep(_pollRate)
	}
}



func Pr(){
	fmt.Println("a")
}


