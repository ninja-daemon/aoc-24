package aocclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func ExtractData(url string) ([][]rune, error) {
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

	return parseInput(string(body)), nil
}

func parseInput(data string) [][]rune {
	lines := strings.Split(strings.TrimSpace(data), "\n")

	matrix := make([][]rune, len(lines))

	for i, line := range lines {
		matrix[i] = []rune(line)
	}

	return matrix
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
