
package FSM

import(
	"fmt"
	"../orders"
	"../elevio"
	"../../timer"
	"../buttons"
	"../../dataTypes"
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"flag"
	"os"
	"./masterCom"
)



var elevator dataTypes.ElevatorInfo

func Init(){
	elevio.SetMotorDirection(dataTypes.D_Down)
	for(elevio.GetFloor() == -1){
		//Wait until reach floor below
		fmt.Println("w")
	}
	elevator.Floor = elevio.GetFloor()
	elevator.CurrentDirection = dataTypes.D_Stop 
	elevator.LocalOrders = [3][4]int{{0,0,0,0},{0,0,0,0},{0,0,0,0}}
	elevator.State = dataTypes.S_Idle

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
	elevator.LocalOrders = orders.ExecuteOrders(elevator)
	elevator.State = dataTypes.S_DoorOpen
	elevator.CurrentDirection = dataTypes.D_Stop
	elevio.SetMotorDirection(elevator.CurrentDirection)
	timer.Start(3000)
	elevio.SetDoorOpenLamp(true)
}

func updateDirection(){
	elevator.CurrentDirection = FindNextDirection(elevator)
	elevio.SetMotorDirection(elevator.CurrentDirection)
	if elevator.CurrentDirection != dataTypes.D_Stop{
		elevator.State = dataTypes.S_Moving
	}
}

func StateMachine(){


	/* ---------- */

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	/* ---------- */












	numFloors := 4
	
    //buttons.Pr()

    elevio.Init("localhost:15657", numFloors)
    
	Init()

    buttonPress           := make(chan dataTypes.ButtonEvent)
	stopButtonPress       := make(chan bool)   
	elevatorInfo          := make(chan dataTypes.ElevatorInfo)
	floorSensor           := make(chan int)
	obstruction           := make(chan bool)
	timedOut              := make(chan bool)

	/* -- Network -- */


	/* ---  -- -- - -*/ 

	go elevio.PollButtons(buttonPress)
	go elevio.PollStopButton(stopButtonPress)
	go buttons.MirrorOrders(elevatorInfo)
	go elevio.PollFloorSensor(floorSensor)
	go elevio.PollObstructionSwitch(obstruction)
	go timer.PollTimer(timedOut)

	/* -- Network -- */


	/* ---  -- -- - -*/ 

	for{
		select{
			case b := <-buttonPress:
				switch b.Button{
				case dataTypes.BT_Cab:
					elevator.LocalOrders[b.Button][b.Floor] = 3
				case dataTypes.BT_HallUp:
					fallthrough
				case dataTypes.BT_HallDown:
					elevator.LocalOrders[b.Button][b.Floor] = 3 // To two later
				}
				/*switch?  */
				if orders.StopHere(elevator) && elevator.State != dataTypes.S_Moving{
					HandleOrder()									
				}else if elevator.State == dataTypes.S_Idle{
					updateDirection()
				}
				
				elevatorInfo <-elevator
				
			case <-timedOut:
				elevio.SetDoorOpenLamp(false)
				updateDirection()
				if elevator.CurrentDirection == dataTypes.D_Stop{
					elevator.State = dataTypes.S_Idle
				}

			case f := <-floorSensor:
				elevator.Floor = f
				elevio.SetFloorIndicator(elevator.Floor)
				fmt.Println("Now on floor")
				if orders.StopHere(elevator){
					HandleOrder()
				}
				elevatorInfo <-elevator

			case <-obstruction:
				dataTypes.ElevatorInfoPrint(elevator)
				
		}
	}

	
}



