
package timer

import "time"


const _pollRate = 20 * time.Millisecond

type door_timer struct{
    start_time_ms int64 
	run_time_ms int64 
	running bool
}

var t = door_timer{}

func Start(run_time_ms int64) {
        
    t.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)

	t.run_time_ms = run_time_ms
	
	t.running = true

}

func TimedOut() bool{


    if ((t.start_time_ms + t.run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}

func PollTimer(receiver chan<- bool){


	for {
		time.Sleep(_pollRate)
		if TimedOut() && t.running{
			receiver <- true
			t.running = false
		}
	}
}


func Hello() string{
    return "world"
}

