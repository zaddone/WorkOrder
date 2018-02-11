package main

import (
	"./request"
	"fmt"
	"time"
)

func main() {
	mr, err := request.HandleOrder(100, request.Instr.MinimumTrailingStopDistance, 0)
	fmt.Println(mr, err)
	time.Sleep(time.Second * 10)
	cr, err := request.ClosePosition()
	fmt.Println(cr, err)
}
