package aocclient

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func ExtractData(url string) ([][]int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	sessionToken, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("error obtaining session token: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading responde: %w", err)
	}

	return parseInput(string(body))
}

func parseInput(data string) ([][]int, error) {
	var result [][]int
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Split(line, " ")

		var row []int

		for _, v := range columns {
			num, err1 := strconv.Atoi(v)
			if err1 != nil {
				return nil, fmt.Errorf("error numeric conversion: %w", err1)
			}
			row = append(row, num)
		}
		result = append(result, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning lines: %w", err)
	}

	return result, nil
}

func getToken() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("file .env not found")
	}
	sessionToken := os.Getenv("SESSION_TOKEN")
	if sessionToken == "" {
		return "", fmt.Errorf("error: session token not available")
	}
	return sessionToken, nil
}
