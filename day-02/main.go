package main

import (
	"day02/aocclient"
	"math"

	"fmt"
	"log"
	"sync"
)

func main() {
	url := "https://adventofcode.com/2024/day/2/input"

	reports, err := aocclient.ExtractData(url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	solution1 := part1Solution(reports)

	fmt.Println("Part 1 solution:", solution1)
}

func part1Solution(reports [][]int) int {
	count := 0
	resultsCh := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := range reports {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resultsCh <- isSafeReport(reports[i])
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		if result {
			count++
		}
	}
	return count
}

func isSafeReport(report []int) bool {
	if len(report) < 2 {
		return false
	}

	isAscending := report[0] < report[1]

	for i := 0; i < len(report)-1; i++ {
		last, curr := report[i], report[i+1]

		if isAscending {
			if curr < last || math.Abs(float64(curr-last)) > 3 || curr == last {
				return false
			}
		} else {
			if curr > last || math.Abs(float64(curr-last)) > 3 || curr == last {
				return false
			}
		}
	}

	return true
}
