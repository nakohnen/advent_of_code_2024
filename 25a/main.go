package main

import (
	"bufio"
	"fmt"
	"os"
    //"sync"
    //"strings"
    "strconv"
)

func readInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		//fmt.Printf("Error converting str %s to int\n", s)
		os.Exit(1)
	}
	return val
}


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

    templatesRaw := [][]string{}
    current := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
        if len(line) > 0 {
            current = append(current, line)
        } else if len(line) == 0 {
            templatesRaw = append(templatesRaw, current)
            current = []string{}
        }
	}

    sum := 0

    keys := [][5]int{}
    locks := [][5]int{}

    for _, template := range templatesRaw {
        isKey := true
        if template[0][0] == '#' {
            isKey = false
        }

        element := [5]int{-1,-1,-1,-1,-1}
        for _, row := range template {
            for i, r := range row {
                if r == '#' {
                    element[i] += 1
                }
            }
        }
        text := ""
        if isKey {
            keys = append(keys, element)
            text = "Key:"
        } else {
            locks = append(locks, element)
            text = "Lock:"
        }
        fmt.Printf("%s %v\n", text, element)

    }
    
    for _, key := range keys {
        for _, lock := range locks {
            fits := true
            for i:=0;i<len(key);i++ {
                if key[i] + lock[i] >= 6 {
                    fits = false
                    break
                }
            }
            if fits {
                sum += 1
            }
        }
    }
   
	fmt.Printf(" -> Sum: %d\n", sum)

    if err4 := scanner.Err(); err4 != nil {
        fmt.Printf("Error reading input file: %v\n", err4)
        os.Exit(1)
    }

	fmt.Println("Processing complete!")
}