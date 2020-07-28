package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var stack []int
var stackSize = 10

//Returns all the cores int the system
func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}
func main() {
	fmt.Printf("GOMAXPROCS is %d\n", getGOMAXPROCS())
	runtime.GOMAXPROCS(1)
	fmt.Printf("GOMAXPROCS is %d", getGOMAXPROCS())
	fmt.Println("\n")
	//Runs with one core
	stack1(1)
	//Runs with multiple cores
	stack1(4)
}

func stack1(numOfCore int) {
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
	fmt.Println("With ", numOfCore, " Cores : ", end.Sub(start), "\n")
}

//Returns true if stack is full
func isFull() bool {
	if len(stack) == stackSize {
		return true
	}
	return false
}

//Returns true if stack is empty
func isEmpty() bool {
	if len(stack) == 0 {
		return true
	}
	return false
}

//Inserts into the Stack
func push(value int) {
	if !isFull() {
		stack = append(stack, value)
		fmt.Printf("Pushed: %d\n", value)
	} else {
		fmt.Printf("Stack is full. Cannot push: %d\n", value)
	}
}

//Deletes the top element of thee stack
func pop() {
	if !isEmpty() {
		top := len(stack) - 1
		fmt.Printf("Popping: %d\n", stack[top])
		stack = stack[:top]
	} else {
		fmt.Printf("Stack is empty. Cannot pop\n")
	}
}

func producer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 6; i++ {
		<-tick
		push(i)
		print()
		ch <- 1
		time.Sleep(1 * time.Second)
	}
	//CLosing the channel
	close(ch)
	wg.Done()
}

func consumer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 6; i++ {
		<-tick
		<-ch
		pop()
		print()
		time.Sleep(2 * time.Second)
	}
	wg.Done()
}

//Print Function
func print() {
	fmt.Println("\n******************** Stack ********************")
	for i := len(stack) - 1; i >= 0; i = i - 1 {
		fmt.Println(" -----")
		if i == len(stack)-1 {
			fmt.Println("|", stack[len(stack)-1], " |", " <-- TOP ")
		} else {
			fmt.Println("|", stack[i], " |")
		}
	}
	fmt.Println(" -----")
	fmt.Println("***********************************************\n")
}
