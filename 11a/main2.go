package main

import (
	"bufio"
	"fmt"
	"os"
	//"sync"
	"strconv"
	"strings"
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

	fmt.Printf("Stones: %v\n", stones)
    //cache := make(map[int][]int)
	for j := 0; j < 75; j++ {

        len_stones := len(stones)
        for i:=0;i<len_stones;i++{
            stone := stones[i]
			if stone == 0 {
                stones[i] = 1
				continue
			}
			digits := fmt.Sprintf("%d", stone)
			if len(digits)%2 == 0 {
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
                stones[i] = val1
				stones = append(stones, val2)
				continue
			}
            stones[i] = stone * 2024
		}
		if len(stones) < 22 {
			fmt.Printf("%d - Stones: %v\n", j, stones)
		} else {

			fmt.Printf("%d - Stones count: %v\n", j, len(stones))
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
