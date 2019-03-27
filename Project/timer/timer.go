
package timer

import "time"


const _pollRate = 20 * time.Millisecond

type Timer_s struct{
    Start_time_ms int 
	Run_time_ms int 
	Running bool
}

func (t *Timer_s) Start(Run_time_ms int) {
	t.Start_time_ms = int(time.Now().UnixNano()) / int(time.Millisecond)
	t.Run_time_ms = Run_time_ms
	t.Running = true
}

func (t* Timer_s) TimedOut() bool{

    if int(t.Start_time_ms) + t.Run_time_ms < int(time.Now().UnixNano()) / int(time.Millisecond) {
        return true
    }
    return false
}

func (t* Timer_s) PollTimer(receiver chan<- bool){

	for {
		time.Sleep(_pollRate)
		if t.TimedOut() && t.Running{
			receiver <- true
			t.Running = false
		}
	}
}

func (t* Timer_s) TimeLeft() int{

	return t.Start_time_ms + t.Run_time_ms - (int(time.Now().UnixNano()) / int(time.Millisecond))
}


func (t* Timer_s) Reset(){
	t.Start_time_ms = int(time.Now().UnixNano()) / int(time.Millisecond)	
	t.Running = true
}

func (t* Timer_s) ResetStartTime(){
	t.Start_time_ms = int(time.Now().UnixNano()) / int(time.Millisecond)	
}