package main

import "gatekeeper/models"

func init() {
	models.Config.Load()
}

func main() {
	for _, listener := range models.Config.Listeners {
		go func(l models.Listener) {
			l.Serve()
		}(listener)
	}

	waitChan := make(chan struct{})
	<-waitChan
}
