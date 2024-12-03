package main

import (
	"day03/aocclient"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {
	url := "https://adventofcode.com/2024/day/3/input"

	memory, err := aocclient.ExtractData(url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	solution1 := part1Solution(memory)

	fmt.Println("Part 1 solution:", solution1)

	solution2 := part2Solution(memory)

	fmt.Println("Part 2 solution:", solution2)

}

func part1Solution(memory string) int {
	re := regexp.MustCompile(`mul\(\d+,\d+\)`)
	matches := re.FindAllString(memory, -1)
	if matches == nil {
		log.Println("Did not find matches")
		return 0
	}

	count := 0
	resultsCh := make(chan int)
	wg := &sync.WaitGroup{}

	for i := range matches {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result, _ := evaluateMulExpression(matches[i])
			resultsCh <- result
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

func part2Solution(memory string) int {
	re := regexp.MustCompile(`mul\(\d+,\d+\)|do\(\)|don\'t\(\)`)
	matches := re.FindAllString(memory, -1)
	if matches == nil {
		log.Println("Did not find matches")
		return 0
	}

	count := 0

	resultsCh := make(chan int)
	wg := &sync.WaitGroup{}

	process := true

	for i := range matches {
		if matches[i] == "do()" {
			process = true
		} else if matches[i] == "don't()" {
			process = false
		}
		if process {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				result, _ := evaluateMulExpression(matches[i])
				resultsCh <- result
			}(i)
		}
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

func evaluateMulExpression(expression string) (int, error) {
	trimmed := strings.TrimSuffix(strings.TrimPrefix(expression, "mul("), ")")
	parts := strings.Split(trimmed, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid expression: %s", expression)
	}

	num1, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, err
	}

	num2, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, err
	}

	return num1 * num2, nil
}
