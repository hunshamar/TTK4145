
package FSM

import(
	"fmt"
	"time"
	"../orders"
	"../elevio"
	"../../timer"
	"../buttonLights"
	"../../config"
	"../../dataTypes"
	"../../network/bcast"
	"../elevatorLogic"
)

const _numFloors int       = config.NumFloors 
const _numElevators int    = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons

var doorTimer     = timer.Timer_s{}
var hardwareTimer = timer.Timer_s{}


func StateMachine(elevatorNumber int, port string){
	elevator := dataTypes.ElevatorInfo{}	
	elevio.Init("localhost:" +  port, _numFloors)
	elevator = Init(elevatorNumber,elevator)

    buttonPress           := make(chan dataTypes.ButtonEvent)
	stopButtonPress       := make(chan bool)   
	floorSensor           := make(chan int)
	obstruction           := make(chan bool)
	timedOut              := make(chan bool)
	TXMaster 			  := make(chan dataTypes.ElevatorInfo)
	RXMaster 			  := make(chan [_numElevators] dataTypes.ElevatorInfo)
	hardwareTimerOut      := make(chan bool)

	go elevio.PollButtons(buttonPress)
	go elevio.PollStopButton(stopButtonPress)
	go elevio.PollFloorSensor(floorSensor)
	go elevio.PollObstructionSwitch(obstruction)
	go doorTimer.PollTimer(timedOut)
	go hardwareTimer.PollTimer(hardwareTimerOut)
	go bcast.Transmitter(config.ElevatorTXPort + elevator.Number, TXMaster)
	go bcast.Receiver(config.MasterTXPort, RXMaster)


	go func(){ // Print elevator status each second
		for{
		time.Sleep(1 * time.Second)
		}
	 }()

	go func(){ // Transmit elevator info to master. 10 hz
		for{
			TXMaster <- elevator
			dataTypes.OrdersPrint(elevator)
			time.Sleep(100 * time.Millisecond)
			buttonLights.MirrorOrders(elevator)
		}
	}()

	
	for{
		select{

		case b := <-buttonPress:
			elevator.LocalOrders = orders.Add(elevator, b.Button, b.Floor)	
			fmt.Println("Button pressed in elevator", elevatorNumber-1)
			if elevatorLogic.ShouldStopHere(elevator) && elevator.State != dataTypes.S_Moving{ // Order on current floor
				elevator = handleOrder(elevator)										
			}else if elevator.State == dataTypes.S_Idle{
				elevator = updateDirection(elevator) 
			}
			
		case f := <-floorSensor:
			elevator.Floor = f
			elevator.HardwareFunctioning = true 
			hardwareTimer.ResetStartTime()
			elevio.SetFloorIndicator(elevator.Floor)

			if elevatorLogic.ShouldStopHere(elevator){
				elevator = handleOrder(elevator)
			}else if (elevator.Floor == 0 || elevator.Floor == _numFloors-1){
				elevio.SetMotorDirection(dataTypes.D_Stop)
				elevator.State = dataTypes.S_Idle
				elevator.CurrentDirection = dataTypes.D_Stop
			}

		case fromMaster := <-RXMaster:	
			elevatorFromMaster := fromMaster[elevator.Number] 
			elevator.LocalOrders = administerOrdersFromMaster(elevator, elevatorFromMaster.LocalOrders) 

			if elevatorLogic.ShouldStopHere(elevator) && elevator.State != dataTypes.S_Moving{
				elevator = handleOrder(elevator)									
			}else if elevator.State == dataTypes.S_Idle{
				elevator = updateDirection(elevator) 
			}		
			
		case <-timedOut:
			elevio.SetDoorOpenLamp(false)
			elevator = updateDirection(elevator)
			if elevator.CurrentDirection == dataTypes.D_Stop{
				elevator.State = dataTypes.S_Idle
			}

		case <- hardwareTimerOut:
			if (elevator.State == dataTypes.S_Moving){
				println("HARDWARE ALERT")
				elevator.HardwareFunctioning = false
			}
		
		case <-obstruction:
			dataTypes.ElevatorInfoPrint(elevator)
		}
	}
		
}


func administerOrdersFromMaster(elevator dataTypes.ElevatorInfo, ordersFromMaster [_numOrderButtons][_numFloors]int) [_numOrderButtons][_numFloors]int{
	localOrders := elevator.LocalOrders
	for floor := 0; floor < _numFloors; floor++{
		for button := 0; button < _numOrderButtons-1; button++{

			switch ordersFromMaster[button][floor]{
			case dataTypes.O_Empty:	
				if localOrders[button][floor] != dataTypes.O_Received{
					localOrders[button][floor] = ordersFromMaster[button][floor]
				}
			case dataTypes.O_LightOn:
				fallthrough
			case dataTypes.O_Handle:
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
	elevator.Number = elevatorNumber-1
	elevio.SetMotorDirection(elevator.CurrentDirection)
	return elevator
}


func handleOrder(elevator dataTypes.ElevatorInfo)dataTypes.ElevatorInfo{
	elevator.LocalOrders = orders.Execute(elevator)
	elevio.SetMotorDirection(dataTypes.D_Stop)
	elevator.State = dataTypes.S_DoorOpen
	elevio.SetDoorOpenLamp(true)
	doorTimer.Start(config.DoorOpenTimeMs)
	return elevator
}

func updateDirection(elevator dataTypes.ElevatorInfo)dataTypes.ElevatorInfo{
	elevator.CurrentDirection = elevatorLogic.FindNextDirection(elevator)
	elevio.SetMotorDirection(elevator.CurrentDirection)
	if elevator.CurrentDirection != dataTypes.D_Stop{
		elevator.State = dataTypes.S_Moving
		hardwareTimer.Start(config.ElevatorTravelTimeMs * 2)
	}
	return elevator
}
