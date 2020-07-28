package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

//Size of the Circular Buffer
const size = 5

//var m = sync.RWMutex{}

var c_buff [size]int
var write, read int

var b bool = false

//Returns all the cores in the system
func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}
func main() {
	fmt.Printf("GOMAXPROCS is %d\n", getGOMAXPROCS())
	runtime.GOMAXPROCS(1)
	fmt.Printf("GOMAXPROCS is %d", getGOMAXPROCS())
	fmt.Println("\n")
	//Runs with one core
	CBuff(1)
	write = 0
	read = 0
	b = false
	//Runs with multiple cores
	CBuff(4)
}

func CBuff(numOfCore int) {
	//Setting Number of cores to run
	runtime.GOMAXPROCS(numOfCore)

	start := time.Now()
	//Waitgroups to run all the go routines
	var wg sync.WaitGroup

	ch := make(chan int, 2)

	wg.Add(2)
	go producer(&wg, ch)
	go consumer(&wg, ch)
	wg.Wait()

	end := time.Now()
	fmt.Println("\nWith ", numOfCore, " Cores : ", end.Sub(start), "\n")
}

// Insert into the Circular Buffer
func enqueue(val int) {
	if write == read && b {
		fmt.Println("Write Pointer Waiting")
		time.Sleep(2 * time.Second)
	} else {
		c_buff[write] = val
		write = (write + 1) % 5
		fmt.Printf("\nQueued value(Producer): %d \n", val)
	}
}

// Read and Delete from the Circular Buffer
func dequeue() {
	if write == read && b {
		fmt.Println("Read Pointer Waiting")
		time.Sleep(2 * time.Second)
	} else {
		val := c_buff[read%5]
		c_buff[read%5] = 0
		fmt.Printf("\nDequeued value(Consumer): %v\n", val)
		read = (read + 1) % 5
		b = true
	}

}

//Producer calls the Enqueue function and sends signal to the channel
func producer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 1; i <= size; i++ {
		<-tick
		enqueue(i)
		//Printing the buffer
		fmt.Println()
		for i := 0; i < size; i++ {
			fmt.Printf("%d ", c_buff[i])
		}
		fmt.Println()
		ch <- 1
		//time.Sleep(2 * time.Second)
	}
	close(ch)
	wg.Done()
}

//Consumer recives the signal from the channel and calls the Dequeue function
func consumer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 1; i <= size; i++ {
		<-tick
		<-ch
		time.Sleep(2 * time.Second)
		dequeue()
		//Printing the Stack
		fmt.Println()
		for i := 0; i < size; i++ {
			fmt.Printf("%d ", c_buff[i])
		}
		fmt.Println()
	}
	wg.Done()
}
