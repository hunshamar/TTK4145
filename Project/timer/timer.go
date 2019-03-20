
package timer

import "time"
//import "fmt"


const _pollRate = 20 * time.Millisecond

type Timer_s struct{
    Start_time_ms int64 
	Run_time_ms int64 
	Running bool
}


func Start(Run_time_ms int64, t* Timer_s) {

	t.Start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)

	t.Run_time_ms = Run_time_ms
	
	t.Running = true

}

func TimedOut(t* Timer_s) bool{


    if ((t.Start_time_ms + t.Run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}

func PollTimer(receiver chan<- bool, t* Timer_s){


	for {
		time.Sleep(_pollRate)
		if TimedOut(t) && t.Running{
			receiver <- true
			t.Running = false
		}
	}
}

func TimedLeft(t Timer_s) int64{

	//fmt.Println("start time", t.Start_time_ms)
	//fmt.Println("Run time", t.Run_time_ms)
	//fmt.Println("Now time", (time.Now().UnixNano() / int64(time.Millisecond)))

	return (t.Start_time_ms + t.Run_time_ms - (time.Now().UnixNano() / int64(time.Millisecond)))
}


func WatchdogTimedOut(w* Timer_s) bool{
    if ((w.Start_time_ms + w.Run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}


func WatchdogInit(time_ms int64, w* Timer_s){
	w.Start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)
	w.Run_time_ms = time_ms
}

func WatchdogPoll(receiver chan<- bool, w* Timer_s){
	for {
		time.Sleep(500*time.Millisecond)
		v := WatchdogTimedOut(w)
		if v == true {
			receiver <- true
		}
	}
}

func WatchdogReset(w* Timer_s){
	w.Start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)	
}

