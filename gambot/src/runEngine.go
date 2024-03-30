package src


import "gambot/src/uci"


func Main() {
	uci.InitEngine()

	stop := false

	//main loop
	for !stop {
		stop = uci.RecieveCmd()  //get and interpret UCI command
	}
}