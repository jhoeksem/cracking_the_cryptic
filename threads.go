package main

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"log"
	//"net/http"
	//"math/rand"
	"sync"
)
func main() {

	var wg sync.WaitGroup

	x := [5]int{10, 20, 30, 40, 50}
	y := [5]int{0, 0, 0, 0, 0}
	y_address := &y

	for i, num := range x {
		fmt.Println(i)
		wg.Add(1)
		go func(num int, index int){
			defer wg.Done()
			(*y_address)[index] = num*num
		}(num, i)
	}
	wg.Wait()
	fmt.Println(y)
}
