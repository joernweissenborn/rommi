package main

import (
	"os"
	"os/signal"
	"rommi/modules/brain/core"
)

func main() {

	if err := core.Start(); err != nil {
		os.Exit(1)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
