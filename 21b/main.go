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


// StringSet type using map[string]struct{}
type StringSet struct {
	m map[string]struct{}
}

// NewStringSet creates a new set
func NewStringSet() *StringSet {
	return &StringSet{m: make(map[string]struct{})}
}

// Add adds an element to the set
func (s *StringSet) Add(value string) {
	s.m[value] = struct{}{}
}

// Contains checks if an element is in the set
func (s *StringSet) Contains(value string) bool {
	_, exists := s.m[value]
	return exists
}

// Remove removes an element from the set
func (s *StringSet) Remove(value string) {
	delete(s.m, value)
}

// Size returns the number of elements in the set
func (s *StringSet) Size() int {
	return len(s.m)
}

// GetElements retrieves all elements in the set
func (s *StringSet) GetElements() []string {
	elements := []string{}
	for p := range s.m {
		elements = append(elements, p)
	}
	return elements
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

// ManhattanDistance calculates the Manhattan distance between two points
func ManhattanDistance(p1, p2 point) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

/**func getPriority(line string, inverse map[string]point, leadingChar string) int {
    newLine := leadingChar + line
    total := 0
    for i:=0;i<len(newLine)-1;i++{
        p1 := inverse[string(newLine[i])]
        p2 := inverse[string(newLine[i+1])]
        total += ManhattanDistance(p1, p2)
    }
    return len(line) + total
}**/

func getPriority(line string, inverse map[string]point, leadingChar string) int {
    if len(line) == 0 {
        return 0
    } else if len(line) == 1 {
        switch line {
        case "<":
            return 0
        case "v":
            return 1
        case ">":
            return 2
        case "^":
            return 3
        case "A":
            return 4
        }
    } 
    return 10 * getPriority(string(line[:len(line)-1]), inverse, leadingChar) + getPriority(string(line[len(line)-1]), inverse, leadingChar)
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

// Priority map for characters
var priorityMap = map[rune]int{
	'<': 1,
	'v': 2,
	'^': 3,
	'>': 4,
	'A': 5,
}

// CompareStrings compares two strings based on their priority
// with a preference for directly repeating characters.
func CompareStrings(s1, s2 string) int {
	// Preference for direct repetitions: Count consecutive repeating characters
	reps1 := countDirectRepetitions(s1)
	reps2 := countDirectRepetitions(s2)
	if reps1 != reps2 {
		return reps2 - reps1 // Prefer the string with more direct repetitions
	}

	// General case: Compare character by character based on priority
	runes1 := []rune(s1)
	runes2 := []rune(s2)

	for i := 0; i < len(runes1) && i < len(runes2); i++ {
		priority1 := priorityMap[runes1[i]]
		priority2 := priorityMap[runes2[i]]

		if priority1 != priority2 {
			return priority1 - priority2 // Negative if s1 < s2, positive if s1 > s2
		}
	}

	// If strings are identical up to the shortest length, return 0 (equal preference)
	return 0
}

// countDirectRepetitions counts the number of directly repeating characters in the string
func countDirectRepetitions(s string) int {
	count := 0
	runes := []rune(s)

	// Check each pair of adjacent characters
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] == runes[i+1] {
			count++
		}
	}

	return count
}

// SortStrings sorts a slice of strings based on the priority and repetition preference
func SortStrings(strings []string) {
	sort.Slice(strings, func(i, j int) bool {
		return CompareStrings(strings[i], strings[j]) < 0
	})
}


/**func filter(perms []string, inverseMap map[string]point, leadingChar string) []string {
    if len(perms) == 0 {
        return perms
    }

    newPerms := []string{}

    minPrio := getPriority(perms[0], inverseMap, leadingChar)
    for _, other := range perms {
        prio := getPriority(other, inverseMap, leadingChar)
        if prio < minPrio {
            minPrio = prio
        }
    }

    for _, other := range perms {

        prio := getPriority(other, inverseMap, leadingChar)
        if prio == minPrio {
            newPerms = append(newPerms, other)
        }
    }
    return newPerms
}**/


func decode(text string, decoder map[string]map[string][]string, inverse map[string]point) []string {
    currentLetter := "A"
    possibilities := []string{}
    for _, r := range text {
        letter := string(r)
        newPoss := decoder[currentLetter][letter]
        if len(possibilities) == 0 {
            possibilities = filter(newPoss, inverse, "A")
        } else {
            toReplace := []string{}
            for _, perm := range possibilities {
                for _, perm2 := range filter(newPoss, inverse, string(perm[len(perm)-1])) {
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

func solve(line string, keypad, dpad map[string]map[string][]string, inverse map[string]point) int {
    ways := decode(line, keypad, inverse)
    newWays := []string{}
    for _, way := range ways {
        newWays = append(newWays, decode(way, dpad, inverse)...)
    }

    minLen := len(newWays[0])
    for _, w := range newWays {
        if len(w) < minLen {
            minLen = len(w)
        }
    }

    return minLen * readInt(line[:len(line)-1])
}

func solve2(line string, keypad map[string]map[string][]string, lenDPad map[string]map[string]int, inverse map[string]point) int {
    ways := decode(line, keypad, inverse)
    
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
                keypad[letter][other] = filter(getPermutations(start, end, invalidKey), layoutKPInverse, "A")
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
                perms := filter(getPermutations(start, end, invalidKey), layoutDPInverse, "A")
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
        for key2, val2 := range val {
            fmt.Printf("%v -> %v : %v\n", key, key2, val2)
        }
    }
    
    // Multi Level DPad transformation
    transform := func(dpadValues []string, dpad, multiLevelDPad map[string]map[string][]string, inverse map[string]point) {
        for _, letter := range dpadValues {
            for _, other := range dpadValues {
                fmt.Printf("%v to %v -> \n", letter, other)
                newPermsSet := NewStringSet()
                for _, currentPerm := range multiLevelDPad[letter][other] {
                    for _, newPerm := range decode(currentPerm, dpad, layoutDPInverse) {
                        newPermsSet.Add(newPerm)
                        fmt.Printf("\t%v to %v\n", currentPerm, len(newPerm))
                    }
                }
                fmt.Printf("New perms size: %d\n", newPermsSet.Size())
                newPerms := newPermsSet.GetElements()
                minPrio := getPriority(newPerms[0], layoutDPInverse, "A")
                minLen := len(newPerms[0])
                for _, perm := range newPerms {
                    prio := getPriority(perm, layoutDPInverse, "A")
                    length := len(perm)
                    if prio < minPrio {
                        minPrio = prio
                    }
                    if length < minLen {
                        minLen = length
                    }
                }
                newPerms2 := []string{}
                for _, perm := range newPerms {
                    prio := getPriority(perm, layoutDPInverse, "A")
                    if prio == minPrio {
                        newPerms2 = append(newPerms2, perm)
                        } else {
                            fmt.Printf("Out: %v %v -> %v\n", letter, other, perm)
                    }
                }
                multiLevelDPad[letter][other] = newPerms2[:1]
                fmt.Printf("%d perms with minLen %d\n", len(newPerms2), minLen)
            }
        }
        for key, val := range multiLevelDPad {
            if len(val) <= 20 {
                for key2, val2 := range val {
                    fmt.Printf("%v -> %v : %v\n", key, key2, val2)
                }
            }
        }
        fmt.Println("")
    }
    for i:=1;i<25;i++{
        // transform := func(dpadValues []string, dpad, multiLevelDPad map[string]map[string][]string, inverse map[string]point) {
        transform(dpadValues, dpad, multiLevelDPad, layoutDPInverse)
    }

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
			result := solve(job, keypad, multiLevelDPad, layoutDPInverse)
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