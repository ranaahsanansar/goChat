package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/utils"
)

func UserGuard(c *gin.Context) {
	// Extracting token
	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	fmt.Println(authHeader)

	// token is in the format "Bearer <token>"
	token := utils.TakeTokenString(authHeader)

	// Check if the token is active
	active, err := IntrospectToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	fmt.Println(active)

	if !active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	} else {

		c.Next()
	}

}

// IntrospectToken verifies the access token using Keycloak's introspection endpoint
func IntrospectToken(token string) (bool, error) {
	// Set up the URL for the introspection endpoint
	introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
		os.Getenv("KEYCLOAK_HOST"), os.Getenv("REALM_NAME"))

	// Create the request body
	data := fmt.Sprintf("client_id=%s&client_secret=%s&token=%s",
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		token)

	// Make the POST request to the introspection endpoint
	req, err := http.NewRequest("POST", introspectionURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	// Parse the response body
	var responseBody struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return false, err
	}

	return responseBody.Active, nil
}
