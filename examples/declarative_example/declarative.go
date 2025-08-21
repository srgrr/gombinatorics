package main

import (
	"context"
	"fmt"
	"sync"

	g "github.com/srgrr/gombinator"
)

func lenTask(wg *sync.WaitGroup, numChan <-chan int) {
	defer wg.Done()
	for numVal := range numChan {
		fmt.Println(numVal)
	}
}

func helloTask(wg *sync.WaitGroup, stringChan <-chan string) {
	defer wg.Done()
	for stringVal := range stringChan {
		fmt.Println(stringVal)
	}
}

func main() {
	ctx := context.Background()
	data := []string{"hello", "world"}
	// Declare channeling functions
	lenChan := g.Map(ctx, data, func(s string) int { return len(s) })
	filterChan := g.Filter(ctx, data, func(s string) bool { return s == "hello" })

	// Run them!
	var wg sync.WaitGroup
	wg.Add(2)
	go lenTask(&wg, lenChan)
	go helloTask(&wg, filterChan)
	wg.Wait()
}
