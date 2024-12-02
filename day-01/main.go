package main

import (
	"day01/aocclient"

	"fmt"
	"log"
	"math"
	"sort"
	"sync"
)

func main() {
	url := "https://adventofcode.com/2024/day/1/input"

	left, right, err := aocclient.ExtractData(url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	sort.Ints(left)
	sort.Ints(right)

	part1Ch := make(chan int)
	part2Ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		part1Ch <- part1Solution(left, right)
	}()

	go func() {
		defer wg.Done()
		part2Ch <- part2Solution(left, right)
	}()

	go func() {
		wg.Wait()
		close(part1Ch)
		close(part2Ch)
	}()

	solution1 := <-part1Ch
	solution2 := <-part2Ch

	fmt.Println("Part 1 solution:", solution1)
	fmt.Println("Part 2 solution:", solution2)
}

func part1Solution(left, right []int) int {
	count := 0
	for i := range left {
		if left[i] != right[i] {
			count += int(math.Abs(float64(left[i] - right[i])))
		}
	}
	return count
}

func part2Solution(left, right []int) int {
	count := 0
	resultsCh := make(chan int)
	wg := &sync.WaitGroup{}

	for i := range left {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resultsCh <- left[i] * similarityCount(left[i], right)
		}(i)
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

func similarityCount(left int, right []int) int {
	var count = 0

	for _, val := range right {
		if left == val {
			count++
		}
	}

	return count
}
