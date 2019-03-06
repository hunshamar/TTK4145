package masterCom

import(
	/*"../bcast"
	"../localip"
	"../peers"
	*/
	"../../dataTypes"
	"time"
)

func Transmit(TX chan<- dataTypes.ShortMessage, elevInfo <-chan dataTypes.ElevatorInfo){

		for {
			HelloMsg := dataTypes.ShortMessage{<-elevInfo}
			TX <- HelloMsg
			time.Sleep(1 * time.Second)
		}
}

