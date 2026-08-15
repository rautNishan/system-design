package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go counter(&wg, "A")
	go counter(&wg, "B")
	wg.Wait()

}

func counter(wg *sync.WaitGroup, name string) {
	for i := 0; i < 1000000; i++ {
		fmt.Printf("Name: %s and %d\n", name, i)
	}
	defer wg.Done()
}
