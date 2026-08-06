package main

import (
	"flag"
	"fmt"
)

var sizeFlag = flag.Int("size", 10, "the size to reduce to")

func main() {
	flag.Parse()
	// args := os.Args[1:]
	userArgs := flag.Args()

	fmt.Println("Hello bitch")
	if len(userArgs) >= 1 {
		fmt.Println("I see you want to open ", userArgs[0])
	}
	fmt.Println("With the size ", *sizeFlag)
	fmt.Println("Duration is", fmt.Sprintf("%.2f", get_duration(userArgs[0]))+"s")
	fmt.Println("Target bitrate is", get_target_bitrate(userArgs[0], float64(*sizeFlag)))
}
