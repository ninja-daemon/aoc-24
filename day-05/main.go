package main

import (
	"day05/aocclient"
	"fmt"
	"log"
	"sync"
)

func main() {
	url := "https://adventofcode.com/2024/day/5/input"

	rules, sequences, err := aocclient.ExtractData(url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	solution1 := part1Solution(rules, sequences)

	fmt.Println("Part 1 solution:", solution1)

	solution2 := part2Solution(rules, sequences)

	fmt.Println("Part 2 solution:", solution2)

}

func part1Solution(rules, sequences [][]int) int {
	count := 0
	resultsCh := make(chan int)
	wg := &sync.WaitGroup{}

	for _, seq := range sequences {
		wg.Add(1)
		go func(seq []int) {
			defer wg.Done()
			if validateSequence(seq, rules) {
				resultsCh <- seq[len(seq)/2]
			}
		}(seq)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		count += result
	}
	return count
}

func part2Solution(rules, sequences [][]int) int {
	count := 0
	resultsCh := make(chan int)
	wg := &sync.WaitGroup{}

	for _, seq := range sequences {
		wg.Add(1)
		go func(seq []int) {
			defer wg.Done()
			if !validateSequence(seq, rules) {
				seq = reorderSequence(seq, rules)
				resultsCh <- seq[len(seq)/2]
			}
		}(seq)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		count += result
	}
	return count
}

func reorderSequence(sequence []int, rules [][]int) []int {
	position := make(map[int]int)
	for i, num := range sequence {
		position[num] = i
	}

	result := append([]int(nil), sequence...)

	for !validateSequence(result, rules) {
		for _, rule := range rules {
			first, second := rule[0], rule[1]
			posFirst, okFirst := position[first]
			posSecond, okSecond := position[second]

			if okFirst && okSecond && posFirst > posSecond {
				position[first], position[second] = posSecond, posFirst
				result[posFirst], result[posSecond] = result[posSecond], result[posFirst]
				break
			}
		}
	}

	return result
}

func validateSequence(sequence []int, rules [][]int) bool {
	position := make(map[int]int)
	for i, num := range sequence {
		position[num] = i
	}

	for _, rule := range rules {
		first, second := rule[0], rule[1]
		posFirst, okFirst := position[first]
		posSecond, okSecond := position[second]
		if okFirst && okSecond && posFirst > posSecond {
			return false
		}
	}

	return true
}
