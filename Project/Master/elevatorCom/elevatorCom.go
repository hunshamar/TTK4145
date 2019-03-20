package elevatorCom

import(
	"../../network/bcast"
	"../../timer"
	"../../dataTypes"
	"../../config"
	"fmt"
	"time"
)








func Receive(RX chan<- dataTypes.ElevatorInfo, elevNumber int, disconnected chan<- int){
	watchdogTimer := timer.Timer_s{}
	timer.WatchdogInit(int64(2000),&watchdogTimer)
	RxElev 				:= make(chan dataTypes.ElevatorInfo)
	watchdogTimedOut 	:= make(chan bool)
	go bcast.Receiver(config.ElevatorTXPort+elevNumber, RxElev)
	go timer.WatchdogPoll(watchdogTimedOut,&watchdogTimer)

	//a := 0


	go func(){
		for{
		time.Sleep(time.Millisecond*350)	
		//fmt.Println("Elevator", elevNumber, timer.TimedLeft(watchdogTimer), "\n")	
		}
	}()

	for {
		select {
			// Greit at heiser ikke kan ta imot hall bestillinger når hardwarefunctioning = false? Tror lett å fikse
		case FromElevator := <-RxElev:
			/*a++
			fmt.Println("Elevator", elevNumber ,a)
			
			*/
			FromElevator.Number = elevNumber

			if FromElevator.HardwareFunctioning == false{
				FromElevator.State = dataTypes.S_Disconnected
				disconnected <- elevNumber
			}

			RX <- FromElevator
			

			timer.WatchdogReset(&watchdogTimer)

		case <-watchdogTimedOut:
			
			if (elevNumber == 2){
			fmt.Println("Elevator 2 disconnected")
			}
			disconnected <- elevNumber
		}
		
	}



}



