package main

import "gatekeeper/models"

func init() {
	models.Config.Load()
}

func main() {
	for _, listener := range models.Config.Listeners {
		go listener.Serve()
	}

	waitChan := make(chan struct{})
	<-waitChan
}
