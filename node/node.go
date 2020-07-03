package main

import (
	"fmt"
	"imsLb"
	"os"
	"os/signal"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	go imsLb.StartLoadBalanceServer("RoundRobinWithThresholdLimited")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println(s)
}
