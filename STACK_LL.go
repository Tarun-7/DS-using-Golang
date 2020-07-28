package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type stack struct {
	value int
	next  *stack
}

var nodeList *stack

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
	stackll(1)//Runs with multiple cores
	stackll(4)
}

func stackll(numOfCore int) {
	//Setting Number of cores to run
	runtime.GOMAXPROCS(numOfCore)
	start := time.Now()
	//Waitgroups to run all the go routines
	var wg sync.WaitGroup

	ch := make(chan int, 10)

	wg.Add(2)
	go producer(&wg, ch)
	go consumer(&wg, ch)
	wg.Wait()

	//Calculates time since start
	end := time.Now()
	fmt.Println("\nWith ", numOfCore, " Cores : ", end.Sub(start), "\n")
}

// PUSH function - Creates Node and Calls Nodelist

func push(val int) {
	node := &stack{val, nil}
	nodeList = Link(node, nodeList)
	//fmt.Println("Value pushed: ", val)
}

//Adds next node address value to the node
func Link(node, nodeList *stack) *stack {
	if nodeList == nil {
		nodeList = node
		return nodeList
	}

	for p := nodeList; p != nil; p = p.next {
		if p.next == nil {
			p.next = node
			return nodeList
		}
	}
	return nodeList
}

// POP function - Removes element from Stack

func pop() *stack {
	//fmt.Println("Popping")
	var prev *stack
	if nodeList == nil {
		return nodeList
	}
	for p := nodeList; p.next != nil; {
		prev = p
		p = p.next
	}
	prev.next = nil
	return nodeList
}

// Printing Stack
func print() {
	fmt.Println("\n----------------- LINKED LIST -----------------\n")
	for p := nodeList; p != nil; p = p.next {
		fmt.Print("->", p.value)
	}
	fmt.Println("\n\n-----------------------------------------------")
}

//Producer calls the push function and sends signal to the channel
func producer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 7; i++ {
		<-tick
		fmt.Println("\nPushing to Stack(Producer): ", i)
		push(i)
		print()
		ch <- 1
		time.Sleep(1 * time.Second)
	}
	close(ch)
	wg.Done()
}

//Consumer recives the signal from the channel and calls the pop function
func consumer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 6; i++ {
		<-tick
		<-ch
		time.Sleep(2 * time.Second)
		fmt.Println("\nPopping(Consumer)")
		pop()
		print()
		//time.Sleep(2 * time.Second)
	}
	wg.Done()
}
