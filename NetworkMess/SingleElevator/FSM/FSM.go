
package FSM

import "fmt"
import "../orders"
import "../elevio"
import "../../timer"
import "../buttons"
import "../../dataTypes"


var elevator dataTypes.ElevatorInfo

func Init(){
	elevio.SetMotorDirection(dataTypes.D_Down)
	for(elevio.GetFloor() == -1){
		//Wait until reach floor below
	}
	elevator.Floor = elevio.GetFloor()
	elevator.CurrentDirection = dataTypes.D_Stop 
	elevator.LocalOrders = [3][4]int{{0,0,0,0},{0,0,0,0},{0,0,0,0}}
	elevator.State = dataTypes.Idle

	elevio.SetMotorDirection(elevator.CurrentDirection)
	fmt.Println("Init complete")
	dataTypes.ElevatorInfoPrint(elevator)
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


func HandleOrder(){
	timer.Start(3000)
	elevio.SetDoorOpenLamp(true)
	elevio.SetMotorDirection(dataTypes.D_Stop)
	//orders.Remove(current_floor)
}

func StateMachine(){



	numFloors := 4
	
    //buttons.Pr()

    elevio.Init("localhost:15657", numFloors)
    
	Init()

    buttonPress           := make(chan dataTypes.ButtonEvent)
	stopButtonPress       := make(chan bool)   
	elevatorInfo          := make(chan dataTypes.ElevatorInfo)
	orderToHandle         := make(chan bool)

	go elevio.PollButtons(buttonPress)
	go elevio.PollStopButton(stopButtonPress)
	go buttons.MirrorOrders(elevatorInfo)

	for{
		select{
			case b := <-buttonPress:
				switch b.Button{
				case dataTypes.BT_Cab:
					elevator.LocalOrders[b.Button][b.Floor] = 3
				case dataTypes.BT_HallUp:
					fallthrough
				case dataTypes.BT_HallDown:
					elevator.LocalOrders[b.Button][b.Floor] = 2
				}

				elevatorInfo <- elevator
			case <-stopButtonPress:
				dataTypes.ElevatorInfoPrint(elevator)
				
			case <- orderToHandle:
				fmt.Println("handle this")
		}
	}
	
	/*
    drv_buttons := make(chan elevio.ButtonEvent)
    floorSensor  := make(chan int)
    drv_obstr   := make(chan bool)
	drv_stop    := make(chan bool)   
	drv_timer   := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(floorSensor)
    go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go buttons.MirrorOrders()
	go timer.PollTimer(drv_timer)

	Init()

	for {   


	

        select {
		
		case a := <- drv_timer: // timer out
			

        case a := <- drv_buttons: 
           

		case floor := <- floorSensor: // Nytt floor, kan jeg stoppe her?
			


        case a := <- drv_obstr:
			
        case a := <- drv_stop:
			
        }

    }   
	*/
	
}



