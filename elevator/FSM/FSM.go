
package FSM

import "fmt"
import "../orders"
import "../elevio"
import "../timer"
import "../buttons"

type state_t int 
const (
	moving_down = -1
	idle = 0
	moving_up = 1	
	door_open = 2
)

var current_state  state_t = 1
var current_floor = 0;

func Init(){
	elevio.SetMotorDirection(elevio.MD_Down)
	for(elevio.GetFloor() == -1){
		//Wait until reach floor below
	}
	current_floor = elevio.GetFloor()
	elevio.SetMotorDirection(elevio.MD_Stop)
	fmt.Println("Init complete")
}

func stateToString(s state_t) string{
	state := s
	switch state{
	case moving_up:
		return "moving up"
	case moving_down:
		return "moving down"
	case idle:
		return "idle"
	case door_open:
		return "door open"
		
	}
	return "error"
}

func PrintStates(){
	fmt.Println("Current floor:", current_floor)
	fmt.Println("Current state:", stateToString(current_state))
	orders.PrintOrders()
}

func Loop(){

	numFloors := 4

    buttons.Pr()


    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go buttons.MirrorOrders()
	

	Init()

	for {   
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            orders.Add(int(a.Button), int(a.Floor))
            
		case a := <- drv_floors:
			
			current_floor = a
			if orders.Local_orders[elevio.BT_Cab][current_floor] == 1{
				PrintStates()
				timer.Start(3000)
				elevio.SetMotorDirection(elevio.MD_Stop)
				orders.Remove(elevio.BT_Cab,current_floor)
				for (!timer.TimedOut()){
					//Wait until timer finished
				} 
				elevio.SetMotorDirection(d)
			}



            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }

    }   
	
	
}



