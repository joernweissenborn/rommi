package main

import (
	"os"
	"os/signal"
	"rommi/ears/core"
)

func main() {
	path, _ := os.Getwd()
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if err := core.Start(path); err != nil {
		os.Exit(1)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
