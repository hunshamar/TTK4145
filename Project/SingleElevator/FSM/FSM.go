
package FSM

import(
	"fmt"
	"time"
	"../orders"
	"../elevio"
	"../../timer"
	"../buttons"
	"../../config"
	"../../dataTypes"
	"../../network/bcast"
	"../elevatorLogic"
)

const _numFloors int = config.NumFloors 
const _numElevators int = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons



func StateMachine(elevatorNumber int, port string){
	elevator := dataTypes.ElevatorInfo{}	
	watchdogTimer := timer.Timer_s{}

	elevio.Init("localhost:" +  port, _numFloors)
	elevator = Init(elevatorNumber,elevator)
	timer.WatchdogInit(int64(5000),&watchdogTimer)



	//hardWareAlert := false 

    buttonPress           := make(chan dataTypes.ButtonEvent)
	stopButtonPress       := make(chan bool)   
	floorSensor           := make(chan int)
	obstruction           := make(chan bool)
	timedOut              := make(chan bool)
	TXMaster 			  := make(chan dataTypes.ElevatorInfo)
	RXMaster 			  := make(chan dataTypes.AllElevatorInfo)
	watchdogTimedOut   	  := make(chan bool)


	


	go timer.WatchdogPoll(watchdogTimedOut,&watchdogTimer)
	go elevio.PollButtons(buttonPress)
	go elevio.PollStopButton(stopButtonPress)
	go elevio.PollFloorSensor(floorSensor)
	go elevio.PollObstructionSwitch(obstruction)
	go timer.PollTimer(timedOut)


	go bcast.Transmitter(config.ElevatorTXPort + elevatorNumber, TXMaster)
	go bcast.Receiver(config.MasterTXPort, RXMaster)


	go func(){ // Printer heisstatus 1 gang i sekundet
		for{
		dataTypes.OrdersPrint(elevator)
		time.Sleep(1 * time.Second)
		}
	 }()

	go func(){ // Sender info til master 10 ganger i sekundet
		for{
			TXMaster <- elevator
			time.Sleep(100 * time.Millisecond)
			buttons.MirrorOrders(elevator)
		}
	}()

	
	for{
		select{
			case b := <-buttonPress:
				fmt.Println("Button pressed in elevator", elevatorNumber)
				switch b.Button{
				case dataTypes.BT_Cab:
					elevator.LocalOrders[b.Button][b.Floor] = dataTypes.O_Handle
				case dataTypes.BT_HallUp:
					fallthrough
				case dataTypes.BT_HallDown:
					elevator.LocalOrders[b.Button][b.Floor] = dataTypes.O_Received 
				}

				if elevatorLogic.ShouldStopHere(elevator) && elevator.State != dataTypes.S_Moving{ // Order on current floor
					elevator = HandleOrder(elevator)									
				}else if elevator.State == dataTypes.S_Idle{
					elevator = updateDirection(elevator) // Bestilling når idle
				}
				
				
				
			case <-timedOut:
				elevio.SetDoorOpenLamp(false)
				elevator = updateDirection(elevator)
				if elevator.CurrentDirection == dataTypes.D_Stop{
					elevator.State = dataTypes.S_Idle
				}

			case f := <-floorSensor:

				/*
				if hardWareAlert{
					fmt.Println("Hardware alert")
					hardWareAlert = false
					elevator.HardwareFunctioning = true
					elevator.State = dataTypes.S_Idle
					elevio.SetMotorDirection(dataTypes.D_Stop)
				
				}*/
				

				elevator.Floor = f
				elevio.SetFloorIndicator(elevator.Floor)
				if elevatorLogic.ShouldStopHere(elevator){
					elevator = HandleOrder(elevator)
				}
				timer.WatchdogReset(&watchdogTimer)


			case <-obstruction:
				dataTypes.ElevatorInfoPrint(elevator)

			 
	
			case fromMaster := <-RXMaster:	
				elevatorFromMaster := fromMaster.Elevators[elevatorNumber -1] 
				elevator.LocalOrders = newOrdersFromMaster(elevator, elevatorFromMaster.LocalOrders) 

				
				if elevatorLogic.ShouldStopHere(elevator) && elevator.State != dataTypes.S_Moving{
					elevator = HandleOrder(elevator)									
				}else if elevator.State == dataTypes.S_Idle{
					elevator = updateDirection(elevator) 
				}		
			case  <-watchdogTimedOut:
				/*
				
				if hardWareAlert{
					fmt.Println("fucked")
					fmt.Println("fucked")
					fmt.Println("fucked")
					elevator.HardwareFunctioning = false
				}

				if elevator.State == dataTypes.S_Moving{
					fmt.Println("Set til 1 pga moving")
					hardWareAlert = true
					timer.WatchdogReset(&watchdogTimer)
				}*/

		}
	}
		
}


func newOrdersFromMaster(elevator dataTypes.ElevatorInfo, ordersFromMaster [_numOrderButtons][_numFloors]int) [_numOrderButtons][_numFloors]int{
	localOrders := elevator.LocalOrders
	for floor := 0; floor < _numFloors; floor++{
		for button := 0; button < _numOrderButtons-1; button++{

			switch ordersFromMaster[button][floor]{
			case 0:	
				if localOrders[button][floor] != dataTypes.O_Received{
					localOrders[button][floor] = ordersFromMaster[button][floor]
				}
			case 1:
				fmt.Println("ERROR")
			case 2:
				fallthrough
			case 3:
				if (localOrders[button][floor] != dataTypes.O_Executed){
					localOrders[button][floor] = ordersFromMaster[button][floor]
				}
			}
		}
		if ordersFromMaster[int(dataTypes.BT_Cab)][floor] == dataTypes.O_Handle{
			localOrders[int(dataTypes.BT_Cab)][floor] = dataTypes.O_Handle
		}
	}
	return localOrders
}


func Init(elevatorNumber int, elevator dataTypes.ElevatorInfo)dataTypes.ElevatorInfo{
	elevio.SetMotorDirection(dataTypes.D_Down)
	for(elevio.GetFloor() == -1){
		//Wait until reach floor below
	}
	elevator.HardwareFunctioning = true
	elevator.Floor = elevio.GetFloor()
	elevator.CurrentDirection = dataTypes.D_Stop 
	elevator.LocalOrders = [_numOrderButtons][_numFloors]int{} // Clear orders
	elevator.State = dataTypes.S_Idle
	elevator.Number = elevatorNumber

	elevio.SetMotorDirection(elevator.CurrentDirection)
	fmt.Println("Init complete")
	dataTypes.ElevatorInfoPrint(elevator)
	return elevator
}


func HandleOrder(elevator dataTypes.ElevatorInfo)dataTypes.ElevatorInfo{
	elevator.LocalOrders = orders.Execute(elevator)
	elevator.State = dataTypes.S_DoorOpen
	elevio.SetMotorDirection(dataTypes.D_Stop)
	timer.Start(config.DoorOpenTimeMs)
	elevio.SetDoorOpenLamp(true)
	return elevator
}

func updateDirection(elevator dataTypes.ElevatorInfo)dataTypes.ElevatorInfo{
	elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
	elevio.SetMotorDirection(elevator.CurrentDirection)
	if elevator.CurrentDirection != dataTypes.D_Stop{
		elevator.State = dataTypes.S_Moving
	}
	return elevator
}
