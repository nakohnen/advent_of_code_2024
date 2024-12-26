package main

import (
	"bufio"
	"fmt"
	"os"
	//"sync"
	"sort"
	"strconv"
	"strings"
)

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

func readInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		//fmt.Printf("Error converting str %s to int\n", s)
		os.Exit(1)
	}
	return val
}

func getAllTriplets(adj map[string][]string, connected map[string]map[string]bool) [][]string {
	triplets := NewStringSet()

	for p1, others := range adj {
		for _, p2 := range others {
			for _, p3 := range adj[p2] {
				if connected[p1][p3] {
					triplet := []string{p1, p2, p3}
					sort.Strings(triplet)
					triplets.Add(strings.Join(triplet, ","))
				}
			}
		}
	}
	result := triplets.GetElements()
	sort.Strings(result)
	finalResult := [][]string{}
	for _, res := range result {
		triplet := strings.Split(res, ",")
		sort.Strings(triplet)
		finalResult = append(finalResult, triplet)
	}
	return finalResult
}

// Generic contains function
func contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func getMaxNetwork(adj map[string][]string, connected map[string]map[string]bool) string {
	triplets := getAllTriplets(adj, connected)
	fmt.Println("Triplets: ")
	for _, triplet := range triplets {
		fmt.Printf("Triplet %v\n", triplet)
	}
	toMerge := make(map[int]int)
    for i, tripletA := range triplets {
        for j, tripletB := range triplets {
            if i == j {
                continue
            }
			nodes := NewStringSet()
			for _, node := range tripletA {
				nodes.Add(node)
			}
			for _, node := range tripletB {
				nodes.Add(node)
			}
			if nodes.Size() != 4 {
				continue
			}
			difference := []string{}
			common := []string{}
			for _, node := range nodes.GetElements() {
				if contains(tripletA, node) && contains(tripletB, node) {
					common = append(common, node)
				} else {
					difference = append(difference, node)
				}
			}
			toLook := []string{}
			for _, node := range common {
				newCandidate := []string{node}
				for _, node2 := range difference {
					newCandidate = append(newCandidate, node2)
				}
				sort.Strings(newCandidate)
				toLook = append(toLook, strings.Join(newCandidate, ","))
			}
			foundCount := 0
			foundOthers := []int{}
            for k, tripletC := range triplets {
                if k == i || k == j {
                    continue
                }
				textRepr := strings.Join(tripletC, ",")
				for _, cand := range toLook {
					if cand == textRepr {
						foundCount += 1
						foundOthers = append(foundOthers, k)
						break
					}
				}
			}
			if foundCount == 2 {
				toMerge[j] = i
				for _, other := range foundOthers {
					toMerge[other] = i
				}
			}
		}
	}

	reverseMerge := make(map[int][]string)
	results := [][]string{}
	for _, to := range toMerge {
		reverseMerge[to] = []string{}
		reverseMerge[to] = append(reverseMerge[to], triplets[to]...)
	}
	for from, to := range toMerge {
		reverseMerge[to] = append(reverseMerge[to], triplets[from]...)
	}
	for source, nodes := range reverseMerge {
		removeDuplicates := NewStringSet()
		for _, node := range nodes {
			removeDuplicates.Add(node)
		}
		reverseMerge[source] = removeDuplicates.GetElements()
		sort.Strings(reverseMerge[source])
		results = append(results, reverseMerge[source])
	}
	fmt.Print("To Merge:\n")
	for _, merge := range results {
		fmt.Printf("%v\n", merge)
	}

	maxLen := 0
	maxId := -1
	for id, merged := range results {
		if len(merged) > maxLen {
			maxLen = len(merged)
			maxId = id
		}
	}

	return strings.Join(results[maxId], ",")
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

	adj := make(map[string][]string)
	connected := make(map[string]map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)

		splits := strings.Split(line, "-")
		p1 := splits[0]
		p2 := splits[1]

		if _, exists := adj[p1]; exists {
			adj[p1] = append(adj[p1], p2)
		} else {
			adj[p1] = []string{p2}
		}

		if _, exists := adj[p2]; exists {
			adj[p2] = append(adj[p2], p1)
		} else {
			adj[p2] = []string{p1}
		}

	}
	for p1, others := range adj {
		connected[p1] = make(map[string]bool)
		for p2, _ := range adj {
			connected[p1][p2] = false
		}
		for _, p2 := range others {
			connected[p1][p2] = true
		}
	}

	//func getMaxNetwork(adj map[string][]string, connected map[string]map[string]bool) string {
	result := getMaxNetwork(adj, connected)

	fmt.Println(result)

	if err4 := scanner.Err(); err4 != nil {
		fmt.Printf("Error reading input file: %v\n", err4)
		os.Exit(1)
	}

	fmt.Println("Processing complete!")
}
