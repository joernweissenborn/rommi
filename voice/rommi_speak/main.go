package main

import (
	"fmt"
	"os"
	"rommi/voice"
)

func main() {
	v, err := voice.New()
	if err != nil {
		panic(err)
	}
	v.Run()
	fmt.Println("Speaking:", os.Args[1])
	v.SpeakAndWait(os.Args[1])
}
