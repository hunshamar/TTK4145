package main

import "./FSM"




func main(){

    go FSM.StateMachine(1)
    go FSM.StateMachine(2)
    go FSM.StateMachine(3)
    

    for{

    }
   
     
}
