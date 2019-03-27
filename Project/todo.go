

package main 

import "fmt"
import "time"


/* TODO
 - Fikse AllElevators.Elevators[e] osv. - GJORt
 - Endre alle steder der bestillingsstatusverdi er tall til dataTypes.O_Handle osv. - GJORT 
 - Endre så konsekvent med number på heis, fra 0-2 i stedet for 1-3  - GJORT
 - Flere verdier som kan defineres i configen -gjort
 - fjerne unødvendig print og kommentering og funksjoner - gjort mesteparten
 - Optimering av costfunksjon, nedprioritert -gjort 
  - mindre select-cases i fms. Gjøre til funksjoner? -gjort 
*/




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