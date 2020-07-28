package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type stack struct {
	value int
	left  *stack
	right *stack
}

var Btree *stack
var c int = 0

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
	stackBtree(1)
	//Runs with multiple cores
	stackBtree(4)
}

func stackBtree(numOfCore int) {
	//Setting Number of cores to run
	runtime.GOMAXPROCS(numOfCore)
	start := time.Now()
	//Waitgroups to run all the go routines
	var wg sync.WaitGroup

	ch := make(chan int, 10)

	wg.Add(2)
	go producer(&wg, ch)
	//wg.Wait()

	//wg.Add(1)
	go consumer(&wg, ch)
	wg.Wait()

	end := time.Now()
	fmt.Println("\nWith ", numOfCore, " Cores : ", end.Sub(start), "\n")
}

// PUSH function - Inserts element into Stack

func push(val int) {
	node := &stack{val, nil, nil}
	Btree = Link(node, Btree)
	print2D(Btree)
}

//Update adress of the next node
func Link(node, nodeList *stack) *stack {
	if Btree == nil {
		Btree = node
		c += 1
		return Btree
	}

	if c%2 == 0 {
		for p := Btree; p != nil; p = p.left {
			if p.left == nil {
				p.left = node
				c += 1
				return Btree
			}
		}
	} else {
		for p := Btree; p != nil; p = p.right {
			if p.right == nil {
				p.right = node
				c += 1
				return Btree
			}
		}
	}
	return Btree
}

//POP function - Remove from stack

func pop() *stack {
	var node *stack
	if Btree == nil {
		return Btree
	}
	// Writing to the Left node
	if c%2 != 0 {
		for p := Btree; p.left != nil; {
			node = p
			p = p.left
		}
		node.left = nil
		c -= 1
		print2D(Btree)
		return Btree
	} else {
		// Writing to the Right node
		for p := Btree; p.right != nil; {
			node = p
			p = p.right
		}
		node.right = nil
		c -= 1
		print2D(Btree)
		return Btree
	}
	print2D(Btree)
	return Btree
}

func print2DUtil(Btree *stack, space int) {
	var count = 5

	// Base case
	if Btree == nil {
		return
	}

	// Increase distance between levels
	space += count

	// Process right child first
	print2DUtil(Btree.right, space)

	// Print current node after space
	// count

	for i := count; i < space; i++ {
		fmt.Printf(" ")
	}

	fmt.Println(Btree.value)

	// Process left child
	print2DUtil(Btree.left, space)
}

// Wrapper over print2DUtil()

func print2D(Btree *stack) {
	fmt.Println()
	fmt.Println("----------Binary tree----------")
	fmt.Println()
	// Pass initial space count as 0
	print2DUtil(Btree, 0)

}

//Producer calls the PUSH function and sends signal to the channel
func producer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 7; i++ {
		<-tick
		fmt.Println("\nPushing to Stack(Producer): ", i)
		push(i)
		//print()
		ch <- 1
		//time.Sleep(1 * time.Second)
	}
	close(ch)
	wg.Done()
}

//Consumer recives the signal from the channel and calls the POP function
func consumer(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)
	for i := 0; i < 6; i++ {
		<-tick
		<-ch
		time.Sleep(2 * time.Second)
		fmt.Println("\nPopping(Consumer)")
		pop()
		//print()
	}
	wg.Done()
}
