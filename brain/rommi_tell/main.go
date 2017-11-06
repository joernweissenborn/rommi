package main

import (
	"fmt"
	"os"
	"rommi/brain"
)

func main() {
	b, err := brain.New()
	if err != nil {
		panic(err)
	}
	b.Run()
	fmt.Println("Telling:", os.Args[1])
	b.TellCommandAndWait(os.Args[1])
}
