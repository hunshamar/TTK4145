
package timer

import "time"


const _pollRate = 20 * time.Millisecond

type timer_s struct{
    start_time_ms int64 
    run_time_ms int64 
}

var t = timer_s{0,0}
var w = timer_s{0,0}

func Start(run_time_ms int64) {
        
    t.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)

    t.run_time_ms = run_time_ms

}

func WatchdogTimedOut() bool{


    if ((w.start_time_ms + w.run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}

/*
func PollTimer(receiver chan<- bool){
	prev := false
	for {
		time.Sleep(_pollRate)
		v := TimedOut()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}*/

func WatchdogInit(time_ms int64){
	w.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)
	w.run_time_ms = time_ms
}

func WatchdogPoll(receiver chan<- bool){



	for {
		time.Sleep(_pollRate)
		v := WatchdogTimedOut()
		if v == true {
			receiver <- true
		}
	}
}

func WatchdogReset(){
	w.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)	
}