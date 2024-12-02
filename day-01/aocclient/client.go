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

func ExtractData(url string) ([]int, []int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	sessionToken, err := getToken()
	if err != nil {
		return nil, nil, fmt.Errorf("error obtaining session token: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading responde: %w", err)
	}

	return parseInput(string(body))
}

func parseInput(data string) ([]int, []int, error) {
	var left, right []int
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Fields(line)
		if len(columns) != 2 {
			continue
		}

		num1, err1 := strconv.Atoi(columns[0])
		num2, err2 := strconv.Atoi(columns[1])
		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("error numeric conversion: %w %w", err1, err2)
		}
		left = append(left, num1)
		right = append(right, num2)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error scanning lines: %w", err)
	}

	return left, right, nil
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
