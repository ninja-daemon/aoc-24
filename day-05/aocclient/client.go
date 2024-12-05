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

func ExtractData(url string) ([][]int, [][]int, error) {
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

	rules, sequences := parseInput(string(body))

	return rules, sequences, nil
}

func parseInput(input string) ([][]int, [][]int) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var rules [][]int
	var sequences [][]int
	currentSection := &rules

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			currentSection = &sequences
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			rule := make([]int, len(parts))
			for i, part := range parts {
				num, _ := strconv.Atoi(part)
				rule[i] = num
			}
			*currentSection = append(*currentSection, rule)
		} else if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")
			sequence := make([]int, len(parts))
			for i, part := range parts {
				num, _ := strconv.Atoi(part)
				sequence[i] = num
			}
			*currentSection = append(*currentSection, sequence)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input", err)
	}

	return rules, sequences

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
