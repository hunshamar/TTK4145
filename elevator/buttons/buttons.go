
package buttons

import "../orders"
import "../elevio"
import "fmt"
import "time"

const _pollRate = 20 * time.Millisecond

func MirrorOrders(){
	for{
		for button_type := 0; button_type <= 2; button_type++{
			for floor := 0; floor <= 3; floor++{
				if orders.Local_orders[2][floor] == 1{
					elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
				} else{
				elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
				}
			}
		}
		time.Sleep(_pollRate)
	}
}

func Pr(){
	fmt.Println("a")
}


