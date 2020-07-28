package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	go a()
	go b()
	time.Sleep(2 * time.Second)
}

func a() {
	fmt.Println("iN A 1")
	runtime.Gosched()
	fmt.Println("iN A 2")
}

func b() {
	fmt.Println("iN B 1")
	runtime.Gosched()
	fmt.Println("iN B 2")
}
