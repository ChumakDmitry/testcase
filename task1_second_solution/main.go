package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counter, length int64
)

func GetResult(x int64) ([]int64, error) {
	counter, length = 0, x

	var (
		wg          sync.WaitGroup
		numArray    = make([]int64, length)
		resultArray = make([]int64, 0, length)
	)

	for i := range numArray {
		numArray[i] = rand.Int63n(10000) - 5000
	}

	for _, v := range numArray {
		wg.Add(1)
		go func() {
			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
			resultArray = append(resultArray, v)
			atomic.AddInt64(&counter, 1)
			wg.Done()
		}()
	}

	wg.Wait()

	return resultArray, nil
}

func CheckProgress() {
	fmt.Printf("Please wait, %d values left\n", length-counter)
	return
}

func main() {
	fmt.Println("Input length array")
	_, err := fmt.Scanf("%d\n", &length)
	if err != nil {
		log.Fatal("Wrong input")
	}

	fmt.Println("Press enter to check progress")

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			reader.ReadString('\n')
			CheckProgress()
		}
	}()

	result, err := GetResult(length)
	if err != nil {
		log.Fatalf("Critical error: %+v", err)
	}

	fmt.Println(result)
}
