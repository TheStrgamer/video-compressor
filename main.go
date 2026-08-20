package main

import (
	"flag"
	"fmt"
	"strings"
)

var sizeFlag = flag.Int("size", 10, "the size to reduce to")

func main() {
	flag.Parse()
	// args := os.Args[1:]
	userArgs := flag.Args()

	if len(userArgs) >= 1 {
		fmt.Println("I see you want to open ", userArgs[0])
	}
	fmt.Println("With the size ", *sizeFlag)
	var name string = userArgs[0]
	if name != "" {
		info, err := GetVideoInfo(name)
		if err == nil {
			info.printInfo()
		}
		out := strings.ReplaceAll(name, ".mp4", "_compressed.mp4")
		if err := CompressToSize(name, out, float64(*sizeFlag)); err != nil {
			fmt.Println("Compression failed:", err)
		}
	}
}
