package aocclient

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ExtractData(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	sessionToken, err := getToken()
	if err != nil {
		return "", fmt.Errorf("error obtaining session token: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading responde: %w", err)
	}

	return string(body), nil
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
