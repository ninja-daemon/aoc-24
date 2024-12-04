package main

import (
	"day04/aocclient"
	"fmt"
	"log"
	"sync"
)

var directions = [8][2]int{
	{-1, 0}, {1, 0},
	{0, -1}, {0, 1},
	{-1, -1}, {-1, 1},
	{1, -1}, {1, 1},
}

var xdirections = [4][2]int{
	{-1, -1}, {-1, 1},
	{1, -1}, {1, 1},
}

func main() {
	url := "https://adventofcode.com/2024/day/4/input"

	board, err := aocclient.ExtractData(url)
	if err != nil {
		log.Fatalf("Error extracting data: %v", err)
	}

	solution1 := part1Solution(board, "XMAS")

	fmt.Println("Part 1 solution:", solution1)

	solution2 := part2Solution(board, "MAS")

	fmt.Println("Part 2 solution:", solution2)

}

func part1Solution(board [][]rune, word string) int {
	count := 0
	resultsCh := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == rune(word[0]) {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					for _, dir := range directions {
						resultsCh <- search(board, word, i, j, 0, dir)
					}
				}(i)

			}
		}
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

func part2Solution(board [][]rune, word string) int {
	count := 0
	resultsCh := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == rune(word[1]) {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					resultsCh <- xSearch(board, word, i, j)
				}(i)

			}
		}
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

func search(board [][]rune, word string, row, col, index int, dir [2]int) bool {
	if index == len(word) {
		return true
	}
	if row < 0 || col < 0 || row >= len(board) || col >= len(board[0]) || board[row][col] != rune(word[index]) {
		return false
	}

	newRow, newCol := row+dir[0], col+dir[1]

	return search(board, word, newRow, newCol, index+1, dir)
}

func xSearch(board [][]rune, word string, row, col int) bool {
	counter := 0
	for _, v := range xdirections {
		dir := [2]int{v[0] * -1, v[1] * -1}
		if search(board, word, row+v[0], col+v[1], 0, dir) {
			counter++
		}
	}

	return counter == 2
}
