

package main 

import "fmt"
import "time"

func main(){
	a := 1
	go goaf(a)
	for{
		a += 1
		time.Sleep(500*time.Millisecond)
		fmt.Println("Is really", a)
	}
}

func goaf(a int){
	for{
		time.Sleep(500*time.Millisecond)

	fmt.Println("int:",a)
	}
}