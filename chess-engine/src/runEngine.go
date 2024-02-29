package src


var stop bool


func Main() {
	initEngine()

	//main loop
	for !stop {
		recieveCmd()  //get and interpret UCI command
	}
}