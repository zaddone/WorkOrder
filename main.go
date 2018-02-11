package main

import (
	"fmt"
	//	"time"
	"./request"
	//	"./cache"
	//	"flag"
)

func main() {
	//	cache.LoadCache()
	//	cache.CheckSensor()
	go request.ReadTemplate()
	var cmd string
	for {
		fmt.Scanf("%s\r", &cmd)
		if cmd == "show" {
			fmt.Println("begin", len(request.TemplagesLib), "---------")
			count := 0
			for _, ste := range request.TemplagesLib {
				le := len(ste.TemplageList)
				count += le
				rate := ste.Winning[0] / float64(le)
				fmt.Println(ste.Key>>1, ste.Winning, rate)
			}
			fmt.Println("end", count)
			//	fmt.Println(cmd)
		}
	}
}
