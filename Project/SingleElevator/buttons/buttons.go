
package buttons

import "../elevio"
import "../../config"
import "../../dataTypes"

const _numFloors int= config.NumFloors
const _numElevators int = config.NumElevators
const _numOrderButtons int = config.NumOrderButtons

func MirrorOrders(elevator dataTypes.ElevatorInfo ){
		for f := 0; f < _numFloors; f++{
			for b := 0; b < _numOrderButtons; b++{
				if elevator.LocalOrders[b][f] >= dataTypes.O_LightOn {
					elevio.SetButtonLamp(dataTypes.ButtonType(b), f, true)
				} else{
				elevio.SetButtonLamp(dataTypes.ButtonType(b), f, false)
				}
			}
		}
		
	}






