package main

import (
	"bufio"
	"fmt"
	"os"
    "sync"
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

type point struct {
    x, y int
}

// GeneratePermutations generates all permutations of the input slice of runes
func GeneratePermutations(input []rune) [][]rune {
	var result [][]rune
	permute(input, 0, &result)
	return result
}

// permute is a helper function for recursive permutation generation
func permute(input []rune, start int, result *[][]rune) {
	if start == len(input) {
		// Append a copy of the current permutation to the result
		temp := make([]rune, len(input))
		copy(temp, input)
		*result = append(*result, temp)
		return
	}

	for i := start; i < len(input); i++ {
		// Swap elements
		input[start], input[i] = input[i], input[start]

		// Recurse on the next element
		permute(input, start+1, result)

		// Backtrack (undo the swap)
		input[start], input[i] = input[i], input[start]
	}
}

func abs(n int) int {
    if n >= 0 {
        return n
    }
    return -n
}


func getPermutations(start, end, invalid point) []string {
    leftRight := end.x - start.x
    upDown := end.y - start.y
    baseSet := []rune{}
    symbol := '>'
    if leftRight < 0 {
        symbol = '<'
    }
    for i:=0;i<abs(leftRight);i++ {
        baseSet = append(baseSet, symbol)
    }


    symbol = 'v'
    if upDown < 0 {
        symbol = '^'
    }
    for i:=0;i<abs(upDown);i++ {
        baseSet = append(baseSet, symbol)
    }

    perms := GeneratePermutations(baseSet)

    validPerms := []string{}
    for _, perm := range perms {
        s := start
        found := false
        for _, r := range perm {
            switch r {
            case '<':
                s.x -= 1
            case '>':
                s.x += 1
            case '^':
                s.y -= 1
            case 'v':
                s.y += 1
            }
            if s == invalid {
                found = true
                break
            }
        }
        if !found {
            exists := false
            result := string(perm) + "A"
            for _, other := range validPerms {
                if other == result {
                    exists = true
                    break
                }
            }
            if !exists {
                validPerms = append(validPerms, result)
            }
        }
    }
    return validPerms
}

func decode(text string, decoder map[string]map[string][]string) []string {
    currentLetter := "A"
    possibilities := []string{}
    for _, r := range text {
        letter := string(r)
        newPoss := decoder[currentLetter][letter]
        if len(possibilities) == 0 {
            possibilities = newPoss
        } else {
            toReplace := []string{}
            for _, perm := range possibilities {
                for _, perm2 := range newPoss {
                    newPerm := perm + perm2
                    toReplace = append(toReplace, newPerm) 
                }
            }
            possibilities = toReplace
        }
        currentLetter = letter
    }
    return possibilities
}

func decodeToLen(text string, decoder map[string]map[string]int) int {
    currentLetter := "A"
    minLen := 0
    for _, r := range text {
        letter := string(r)
        minLen += decoder[currentLetter][letter]
        currentLetter = letter
    }
    return minLen
}

func solve(line string, keypad, dpad map[string]map[string][]string) int {
    ways := decode(line, keypad)
    newWays := []string{}
    for _, way := range ways {
        newWays = append(newWays, decode(way, dpad)...)
    }

    minLen := len(newWays[0])
    for _, w := range newWays {
        if len(w) < minLen {
            minLen = len(w)
        }
    }

    return minLen * readInt(line[:len(line)-1])
}

func solve2(line string, keypad map[string]map[string][]string, lenDPad map[string]map[string]int) int {
    ways := decode(line, keypad)
    
    waysLen := []int{}
    for _, perm := range ways {
        waysLen = append(waysLen, decodeToLen(perm, lenDPad))
    }

    minLen := waysLen[0]
    for _, l := range waysLen {
        if l < minLen {
            minLen = l
        }
    }
    return minLen * readInt(line[:len(line)-1])
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


    keypad := make(map[string]map[string][]string)
    keypadValues := []string{"A"}
    for i:=0;i<=9;i++{
        val := fmt.Sprintf("%d", i)
        keypad[val] = make(map[string][]string)
        keypadValues = append(keypadValues, val)

    }
    keypad["A"] = make(map[string][]string)
    
    layoutKP := make(map[point]string)
    layoutKPInverse := make(map[string]point)
    layoutKP[point{0, 0}] = "7"
    layoutKP[point{1, 0}] = "8"
    layoutKP[point{2, 0}] = "9"

    layoutKP[point{0, 1}] = "4"
    layoutKP[point{1, 1}] = "5"
    layoutKP[point{2, 1}] = "6"

    layoutKP[point{0, 2}] = "1"
    layoutKP[point{1, 2}] = "2"
    layoutKP[point{2, 2}] = "3"

    layoutKP[point{1, 3}] = "0"
    layoutKP[point{2, 3}] = "A"

    invalidKey := point{0, 3}

    for _, letter := range keypadValues {
        for key, val := range layoutKP {
            if val == letter {
                layoutKPInverse[letter] = key
            }
        }
    }
    // getPermutations(start, end, invalid point) []string 
    for _, letter := range keypadValues {
        for _, other := range keypadValues {
            if letter == other {
                keypad[letter][other] = []string{"A"}
            } else {
                start := layoutKPInverse[letter]
                end := layoutKPInverse[other]
                keypad[letter][other] = getPermutations(start, end, invalidKey)
            }
        }
    }

    dpad := make(map[string]map[string][]string)
    multiLevelDPad := make(map[string]map[string][]string)
    dpadValues := []string{"A", "<", ">", "^", "v"}
    for _, d := range dpadValues {
        dpad[d] = make(map[string][]string)
        multiLevelDPad[d] = make(map[string][]string)
    }

    layoutDP := make(map[point]string)
    layoutDPInverse := make(map[string]point)
    layoutDP[point{1, 0}] = "^"
    layoutDP[point{2, 0}] = "A"

    layoutDP[point{0, 1}] = "<"
    layoutDP[point{1, 1}] = "v"
    layoutDP[point{2, 1}] = ">"

    invalidKey = point{0, 0}

    for _, letter := range dpadValues {
        for key, val := range layoutDP {
            if val == letter {
                layoutDPInverse[letter] = key
            }
        }
    }
    // getPermutations(start, end, invalid point) []string 
    lenDPad := make(map[string]map[string]int)
    for _, letter := range dpadValues {
        lenDPad[letter] = make(map[string]int)
        for _, other := range dpadValues {
            if letter == other {
                dpad[letter][other] = []string{"A"}
                multiLevelDPad[letter][other] = []string{"A"}
            } else {
                start := layoutDPInverse[letter]
                end := layoutDPInverse[other]
                perms := getPermutations(start, end, invalidKey)
                dpad[letter][other] = perms
                multiLevelDPad[letter][other] = append([]string(nil), perms...) 
            }
        }
    }

    fmt.Println("Keypad:")
    for key, val := range keypad {
        fmt.Printf("%v -> %v\n", key, val)
    }
    fmt.Println("Dpad:")
    for key, val := range dpad {
        fmt.Printf("%v -> %v\n", key, val)
    }
    
    // Multi Level DPad transformation
    for i:=1;i<25;i++{
        fmt.Printf("Dpad level %d\n", i + 1)
        for _, letter := range dpadValues {
            for _, other := range dpadValues {
                fmt.Printf("%d: %v to %v -> ", i + 1, letter, other)
                newPerms := []string{}
                for _, currentPerm := range multiLevelDPad[letter][other] {
                    newPerms = append(newPerms, decode(currentPerm, dpad)...)
                }
                minLen := len(newPerms[0])
                for _, perm := range newPerms {
                    if len(perm) < minLen {
                        minLen = len(perm)
                    }
                }
                newPerms2 := []string{}
                for _, perm := range newPerms {
                    if len(perm) == minLen {
                        newPerms2 = append(newPerms2, perm)
                    }
                }
                lenDPad[letter][other] = minLen
                multiLevelDPad[letter][other] = newPerms2
                fmt.Printf("%d perms with minLen %d\n", len(newPerms2), minLen)
            }
        }
    }
    fmt.Println("Multi level dpad done")



    // Prepare for line-by-line reading and writing
    scanner := bufio.NewScanner(inFile)

    sum := 0

    toWork := []string{}
    for scanner.Scan() {
        line := scanner.Text()
        toWork = append(toWork, line)
    }
    
	// Worker function closure that uses `multiplier`
	worker := func(id int, jobs <-chan string, results chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()

		for job := range jobs {
			result := solve2(job, keypad, lenDPad)
			fmt.Printf("Worker %d processed job: %s -> %d\n", id, job, result)
			results <- result
		}
	}

	// Inputs and setup
	jobs := make(chan string, len(toWork))
	results := make(chan int, len(toWork))
	var wg sync.WaitGroup

	// Start workers
	numWorkers := 14
    if numWorkers > len(toWork) {
        numWorkers = len(toWork)
    }
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Send jobs
	for _, input := range toWork {
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
