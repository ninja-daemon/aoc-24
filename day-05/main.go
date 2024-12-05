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

func validateSequence(sequence []int, rules [][]int) bool {
	position := make(map[int]int)
	for i, num := range sequence {
		position[num] = i
	}

	for _, rule := range rules {
		first, second := rule[0], rule[1]
		if posFirst, okFirst := position[first]; okFirst {
			if posSecond, okSecond := position[second]; okSecond {
				if posFirst > posSecond {
					return false
				}
			}
		}
	}
	return true
}
