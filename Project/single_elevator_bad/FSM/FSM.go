
package FSM

import "fmt"
import "../orders"
import "../elevio"
import "../timer"
import "../buttons"
import "../dataTypes"


var current_state  dataTypes.State_t = 1
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

func stateToString(s dataTypes.State_t) string{
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

func FindNextState(floor int, current_direction elevio.MotorDirection) dataTypes.State_t{
	if (orders.Empty()){
		return idle
	}
	
	if (current_direction == elevio.MD_Stop && orders.OnFloor(floor)){
		return idle
	}

	if (current_direction == elevio.MD_Up || current_direction == elevio.MD_Stop){
		if orders.Above(floor){
			return moving_up
		}else{
			return moving_down // snu, finnes ordre andre veien
		}
	}


	if (current_direction == elevio.MD_Down){
		if orders.Below(floor){
			return moving_down
		}else{
			return moving_up // snu
		}
	}
	return idle
}

func HandleOrder(){
	timer.Start(3000)
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(elevio.MD_Stop)
	orders.Remove(current_floor)
}

func Loop(){

	numFloors := 4

    buttons.Pr()


    elevio.Init("localhost:15657", numFloors)
    
	var d elevio.MotorDirection = elevio.MD_Stop

    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
	drv_stop    := make(chan bool)   
	drv_timer   := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go buttons.MirrorOrders()
	go timer.PollTimer(drv_timer)

	Init()

	for {   


		if elevio.GetFloor() != -1 && timer.TimedOut(){
			
			current_state = FindNextState(current_floor,d)
			switch current_state{
				case moving_up:
					d = elevio.MD_Up
				case moving_down:
					d = elevio.MD_Down
				case idle:
					d = elevio.MD_Stop
				case door_open:
					d = elevio.MD_Stop
			}
			elevio.SetMotorDirection(d)
		}

        select {
		
		case a := <- drv_timer:
			if (a){
				elevio.SetDoorOpenLamp(false)
			}else{
				elevio.SetDoorOpenLamp(true)
			}

        case a := <- drv_buttons: 
            fmt.Printf("%+v\n", a)
			orders.Add(int(a.Button), int(a.Floor))

			
			if orders.ExecutableOnFloor(int(current_state), current_floor) && elevio.GetFloor() != -1{
				current_state = door_open
				HandleOrder()			
			}
			
			
		

		case a := <- drv_floors: // Nytt floor, kan jeg stoppe her?
			current_floor = a
			elevio.SetFloorIndicator(current_floor)
			

			if orders.ExecutableOnFloor(int(current_state), current_floor){
				current_state = door_open
				HandleOrder()
			}



        case a := <- drv_obstr:
			fmt.Printf("%+v\n", a)
			PrintStates()
			/*
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }*/
            
        /*case a := <- drv_stop:
			

			
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }*/
        }

    }   
	
	
}



