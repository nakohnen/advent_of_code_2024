package main

import (
	"fmt"
)

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

func main() {
	// Example input
	list := []rune{'a', 'b', 'c'}

	// Generate permutations
	permutations := GeneratePermutations(list)

	// Print results
	fmt.Println("Permutations:")
	for _, perm := range permutations {
		fmt.Println(string(perm)) // Convert runes back to string for display
	}
}
