package main

import (
	"bufio"
	"fmt"
	"os"
    "sync"
    "strings"
    "strconv"
    "sort"
)

func readInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		//fmt.Printf("Error converting str %s to int\n", s)
		os.Exit(1)
	}
	return val
}

func isPossible(towels []string, design string) bool {
    possible := false

    toWork := []string{design}

    for len(toWork) >  0 {
        current := toWork[0]
        toWork = toWork[1:]

        for _, t := range towels {
            found := true
            if len(t) > len(current) {
                continue
            }
            for i:=0;i<len(t);i++{
                if t[i] != current[i] {
                    found = false
                }
            }

            if found {
                subDesign := current[len(t):]
                //fmt.Printf("%v -> %v with %v\n", current, subDesign, t)
                if len(subDesign) == 0 {
                    return true
                }
                toWork = append(toWork, subDesign)
            }
        }
    }

    return possible
}

func removeDuplicates(slice [][]string) [][]string {
    result := [][]string{}
    sliceSimple := []string{}
    ignore := make(map[int]bool)
    for i, sub := range slice {
        ignore[i] = false
        comb1 := strings.Join(sub, "_")
        sliceSimple = append(sliceSimple, comb1)
    }


    for i, sub := range sliceSimple {
        for j:=i+1;j<len(sliceSimple);j++{
            if !ignore[j] && sub == sliceSimple[j] {
                ignore[j] = true
            }
        }
    }

    for i, sub := range slice {
        if !ignore[i] {
            result = append(result, sub)
        }
    }
    return result
}

func findCombination(towels []string, design string, removed []string) int {
    toWork := []string{design}
    var toWorkComp [][]string
    toWorkComp = append(toWorkComp, []string{})

    done := [][]string{}
    for len(toWork) >  0 {
        current := toWork[0]
        toWork = toWork[1:]

        comp := toWorkComp[0]
        toWorkComp = toWorkComp[1:]

        for _, t := range towels {
            found := true
            if len(t) > len(current) {
                continue
            }
            for i:=0;i<len(t);i++{
                if t[i] != current[i] {
                    found = false
                }
            }

            if found {
                subDesign := current[len(t):]
                comp = append(comp, t)
                //fmt.Printf("%v -> %v with %v\n", current, subDesign, t)
                if len(subDesign) == 0 {
                    done = append(done, comp)
                } else {
                    toWork = append(toWork, subDesign)
                    toWorkComp = append(toWorkComp, comp)


                }
            }
        }
    }

    done2 := [][]string{}

    for _, comp := range done {
        
        removed2 := []string{}
        for _, rt := range removed {
            if isPossible(comp, rt) {
                removed2 = append(removed2, rt)
            }
        }

        toWork := [][]string{comp} 

        for len(toWork)>0 {
            current := toWork[0]
            toWork = toWork[1:]
            done2 = append(done2, current)

            for i:=0;i<len(current);i++ {
                for j:=i;j<len(current);j++ {
                    combined := strings.Join(current[i:j] , "")

                    for _, rt := range removed2 {
                        if combined == rt {

                            newComp := append([]string(nil), current[:i]...)
                            newComp = append(newComp, rt)
                            newComp = append(newComp, current[j+1:]...)
                            toWork = append(toWork, newComp)
                        }
                    }
                }
            }
        }
    }

    done2 = removeDuplicates(done2)

    return len(done2)
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

	sum := 0

    towels := []string{}

    designs := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
        if len(towels) == 0 {
            towels = strings.Split(line, ", ")
        } else if len(line) > 0 {
            designs = append(designs, line)
        }
	}

	// Sort by length and then alphanumerically
	sort.Slice(towels, func(i, j int) bool {
		if len(towels[i]) == len(towels[j]) {
			return towels[i] < towels[j] // Alphanumerical for equal lengths
		}
		return len(towels[i]) < len(towels[j]) // By length
	})

    fmt.Printf("Towels %v\n", towels)
    //shorten towels list
    removeIndexes := make(map[int]bool)
    for i, t := range towels {
        //fmt.Printf("%v %v\n", i, t)
        removeIndexes[i] = false
        otherTowels := append([]string(nil), towels[:i]...)
        otherTowels = append(otherTowels, towels[i+1:]...)
        //fmt.Printf("Other towels %v with %v\n", otherTowels, t)
        if isPossible(otherTowels, t) {
            removeIndexes[i] = true
        }
    }
    towels2 := []string{}
    removed := []string{}
    for i, towel := range towels {
        if !removeIndexes[i] {
            towels2 = append(towels2, towel)
        } else {
            removed = append(removed, towel)
            fmt.Printf("Removing %s at %d\n", towel, i) 
        }
    }

    fmt.Printf("Towels %v\n", towels2)
    //fmt.Printf("Designs %v\n", designs)

    worker := func(id int, jobs <-chan string, results chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()

		for job := range jobs {
            result := 0
                if isPossible(towels2, job) {
                    result = findCombination(towels2, job, removed)
                }
			fmt.Printf("Worker %d processed job: %v -> %d\n", id, job, result)
			results <- result
		}
	}

	// Inputs and setup
	jobs := make(chan string, len(designs))
	results := make(chan int, len(designs))
	var wg sync.WaitGroup

	// Start workers
	numWorkers := 24
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Send jobs
	for _, input := range designs {
		jobs <- input
	}
	close(jobs)

	// Wait and collect results
	wg.Wait()
	close(results)

	for result := range results {
        sum += result
	}
   
	fmt.Printf(" -> Sum: %d\n", sum)

    if err4 := scanner.Err(); err4 != nil {
        fmt.Printf("Error reading input file: %v\n", err4)
        os.Exit(1)
    }

	fmt.Println("Processing complete!")
}
