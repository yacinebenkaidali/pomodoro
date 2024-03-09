package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	count := 0

	go func() {
		for v := range time.Tick(1 * time.Second) {
			fmt.Printf("%v\n", v.Format("4:05"))
			count++
			if count == 10 {
				done <- struct{}{}
			}
		}
	}()
	<-done

	fmt.Println("10 seconds have passed")
}
