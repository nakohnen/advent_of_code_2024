package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// Check if enough arguments are provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		os.Exit(1)
	}

	// Parse command-line arguments
	inputFile := os.Args[1]

	// Open the input file
	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	// Prepare for line-by-line reading and writing
	scanner := bufio.NewScanner(inFile)

	sum := 0

	stones := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)

		for _, part := range strings.Split(line, " ") {
			val, err := strconv.Atoi(part)
			if err != nil {
				fmt.Printf("Error converting number %v\n", part)
				os.Exit(1)
			}
			stones = append(stones, val)
		}
	}

    cache := make(map[int][]int)
	var rwMutex sync.RWMutex

	fmt.Printf("Stones: %v\n", stones)
	for i := 0; i < 75; i++ {
		// Worker function closure
		worker := func(id int, jobs <-chan int, results chan<- []int, wg *sync.WaitGroup, cache map[int][]int, rwMutex *sync.RWMutex) {
			defer wg.Done()

			for stone := range jobs {
                rwMutex.RLock()
                if cached, found := cache[stone]; found {

                    newStones := []int{}
                    newStones = append(newStones, cached...)
                    rwMutex.RUnlock()
                    results <- newStones
                    continue
                }
                rwMutex.RUnlock()

				newStones := []int{}
				finished := false
				if stone == 0 {
					newStones = append(newStones, 1)
					finished = true
				}
				digits := fmt.Sprintf("%d", stone)
				if len(digits)%2 == 0 && !finished {
					val1, err := strconv.Atoi(digits[:len(digits)/2])
					if err != nil {
						fmt.Printf("Error converting number %s\n", digits)
						os.Exit(1)
					}
					val2, err := strconv.Atoi(digits[len(digits)/2:])
					if err != nil {
						fmt.Printf("Error converting number %s\n", digits)
						os.Exit(1)
						os.Exit(1)
					}
					newStones = append(newStones, val1)
					newStones = append(newStones, val2)
					finished = true
				}
				if !finished {
					newStones = append(newStones, stone*2024)
				}

                rwMutex.Lock()
                cache[stone] = newStones
                rwMutex.Unlock()

				results <- append([]int(nil), newStones...)
			}
		}

		// Inputs and setup
		jobs := make(chan int, len(stones))
		results := make(chan []int, len(stones))
		var wg sync.WaitGroup

		// Start workers
		numWorkers := 14
		for i := 1; i <= numWorkers; i++ {
			wg.Add(1)
			go worker(i, jobs, results, &wg, cache, &rwMutex)
		}

		// Send jobs
		for _, input := range stones {
			jobs <- input
		}
		close(jobs)

		// Wait and collect results
		wg.Wait()
		close(results)

		newStones := []int{}
		for result := range results {
			for _, newStone := range result {
				newStones = append(newStones, newStone)
			}
		}
		stones = newStones
		if len(stones) < 22 {
			fmt.Printf("%d - Stones: %v\n", i, stones)
		} else {

			fmt.Printf("%d - Stones count: %v\n", i, len(stones))
		}
	}

	sum = len(stones)

	fmt.Printf(" -> Sum: %d\n", sum)

	if err4 := scanner.Err(); err4 != nil {
		fmt.Printf("Error reading input file: %v\n", err4)
		os.Exit(1)
	}

	fmt.Println("Processing complete!")
}
