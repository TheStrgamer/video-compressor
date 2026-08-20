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

	if len(userArgs) >= 1 {
		fmt.Println("I see you want to open ", userArgs[0])
	}
	fmt.Println("With the size ", *sizeFlag)
	var name string = userArgs[0]
	if name != "" {
		fmt.Println("Duration is", fmt.Sprintf("%.2f", get_duration(name))+"s")
		fmt.Println("Target bitrate is", get_target_bitrate(name, float64(*sizeFlag)))
		info, err := GetVideoInfo(name)
		if err == nil {
			info.printInfo()
		}
		if err := CompressToSize(name, "C:/Users/konra/Videos/out3.mp4", float64(*sizeFlag)); err != nil {
			fmt.Println("Compression failed:", err)
		} else {
			fmt.Println("Saved to out.mp4")
		}
	}
}
