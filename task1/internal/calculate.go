package internal

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counter, length int64
)

func GetResult(x int64) (int64, error) {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return x, nil
}

func Calculate(x int) []int64 {
	counter, length = 0, int64(x)

	var (
		wg          sync.WaitGroup
		numArray    = make([]int64, length)
		resultArray = make([]int64, 0, length)
	)

	for i := range numArray {
		numArray[i] = rand.Int63n(10000) - 5000
	}

	for i := 0; i < len(numArray); i++ {
		wg.Add(1)
		go func() {
			result, err := GetResult(numArray[counter])
			atomic.AddInt64(&counter, 1)
			if err != nil {
				fmt.Println(err)
			}

			resultArray = append(resultArray, result)
			wg.Done()
		}()
	}

	wg.Wait()

	return resultArray
}

func CheckProgress() string {
	return fmt.Sprintf("Please wait, %d values left\n", length-counter)
}
