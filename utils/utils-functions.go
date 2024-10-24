package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func PostRequestJson(url string, dataBytes *bytes.Buffer, token string) ([]byte, int, error) {
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, dataBytes)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response body: %v", err)
	}

	// Check for non-success status codes
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("API error: %s", string(body))
	}

	// Return the raw response body
	return body, resp.StatusCode, nil

}

func PostRequestUrlEncoded(url string, data map[string]string, token string) ([]byte, error) {
	requestBody := mapToUrlEncoded(data)
	fmt.Println(requestBody)

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Send the request
	// client := &http.Client{}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Check for non-success status codes
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	// Return the raw body
	return body, nil
}

func mapToUrlEncoded(data map[string]string) string {
	// Create a url.Values object
	values := url.Values{}

	// Loop through the map and add the key-value pairs to the url.Values object
	for key, value := range data {
		values.Set(key, value)
	}

	// Encode the map into a URL-encoded string
	return values.Encode()
}

func GetKeyCloakAdminToken() ([]byte, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", os.Getenv("KEYCLOAK_HOST"), os.Getenv("REALM_NAME"))

	fmt.Println(url)
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("CLIENT_ID"),
		"client_secret": os.Getenv("CLIENT_SECRET"),
	}

	body, err := PostRequestUrlEncoded(url, data, "")
	if err != nil {
		fmt.Println(
			"Error Here: ",
			err)
		return nil, err
	}

	return body, nil

}

// Generic function to unmarshal the raw response body into a given type
func UnmarshalResponse[T any](body []byte, result *T) error {
	return json.Unmarshal(body, result)
}

func CheckPasswordHash(password, hash string) bool {
	if password == hash {
		return true
	} else {
		return false

	}
}

func TakeTokenString(token string) string {
	// Ensure the token has the correct "Bearer" format
	parts := strings.Split(token, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {

		fmt.Println("Invalid token format")
		return ""
	}

	return parts[1]
}
