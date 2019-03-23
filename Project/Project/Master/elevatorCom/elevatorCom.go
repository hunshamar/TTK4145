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
	timer.Start(int64(config.WatchdogTimeMs),&watchdogTimer)
	RxElev 				:= make(chan dataTypes.ElevatorInfo)
	watchdogTimedOut 	:= make(chan bool)
	go bcast.Receiver(config.ElevatorTXPort+elevNumber, RxElev)
	go timer.PollTimer(watchdogTimedOut,&watchdogTimer)

	a := 0


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
			a++
			if elevNumber == 1{

				fmt.Println("Elevator", elevNumber ,a)
			}
			
			FromElevator.Number = elevNumber

			if FromElevator.HardwareFunctioning == false{
				FromElevator.State = dataTypes.S_Disconnected
				disconnected <- elevNumber
			}

			RX <- FromElevator
			

			timer.Reset(&watchdogTimer)

		case <-watchdogTimedOut:
			
			if (elevNumber == 2){
			fmt.Println("Elevator 2 disconnected")
			}
			disconnected <- elevNumber
			timer.Reset(&watchdogTimer) // Må testes mer
		}
		
	}



}


