package main

import (
	srvr "webserver/pkg/router"
)

func main() {
	//serve webserver in the background, allow user input to continue
	srvr.ServeTemplatesAndStyles()

}
