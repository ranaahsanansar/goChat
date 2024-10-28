package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	userDto "github.com/ranaahsanansar/gochat/internal/user/dto"
	"github.com/ranaahsanansar/gochat/models"
	"github.com/ranaahsanansar/gochat/utils"
)

func CreateUser(c *gin.Context) {

	// get body from request

	var body userDto.RegisterUserDTO

	// Bind incoming JSON to the DTO
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate password
	if body.Password != body.ConfirmPassword {
		c.JSON(400, gin.H{
			"message": "passwords do not match",
		})
		return
	}

	// Create in DB
	user := models.Users{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	}

	// check in Db
	isExists := initializers.DB.First(&user, "email = ?", body.Email)

	if isExists.RowsAffected > 0 {
		c.JSON(400, gin.H{
			"message": "user already exists",
		})
		return
	}

	// create in keycloak and then create in DB
	_, err := createUserInKeyCloak(&user)

	if err != nil {
		c.JSON(500, gin.H{
			"message": errors.New("failed to create user").Error(),
		})
		return
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to create user",
		})
		return
	}

	// send response
	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
	})
}

func GetUser(c *gin.Context) {
	// extract userid from param
	var userId = c.Param("userid")
	// search on DB

	var user models.Users

	result := initializers.DB.First(&user, userId)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"message": "user not found",
		})
		return
	}

	// send response

	c.JSON(200, gin.H{
		"message": "user found",
		"data": gin.H{
			"user": user},
	})

}

func LogInUser(c *gin.Context) {

	var body userDto.LoginUserDTO
	// Bind incoming JSON to the DTO
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// search on DB
	var userModel models.Users
	result := initializers.DB.First(&userModel, "username = ?", body.Username)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"message": "user not found",
		})
		return
	}

	// check password
	if utils.CheckPasswordHash(body.Password, userModel.Password) == false {
		c.JSON(401, gin.H{
			"message": "invalid credentials",
		})
		return
	}

	var kyecloakUserCredential = map[string]string{
		"username":      userModel.Username,
		"password":      userModel.Password,
		"grant_type":    "password",
		"client_id":     os.Getenv("CLIENT_ID"),
		"client_secret": os.Getenv("CLIENT_SECRET"),
	}

	keycloakResponse, err := utils.PostRequestUrlEncoded(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", os.Getenv("KEYCLOAK_HOST"), os.Getenv("REALM_NAME")), kyecloakUserCredential, "")
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": "failed to Sing in",
		})
		return
	}

	var keycloakLoginResponse utils.KeycloakLoginResponseDto
	err = json.Unmarshal(keycloakResponse, &keycloakLoginResponse)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": "failed to Sing in",
		})
		return
	}

	// send response
	c.JSON(200, gin.H{
		"message": "user found",
		"data":    keycloakLoginResponse,
	})
}

func RefreshToken(c *gin.Context) {
	// Define your Keycloak token endpoint
	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", os.Getenv("KEYCLOAK_HOST"), os.Getenv("REALM_NAME"))

	refreshToken := c.PostForm("refresh_token")
	frontClientSecret := c.PostForm("client_secret")

	if frontClientSecret != os.Getenv("FRONT_END_CLIENT_SECRET") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client secret"})
		c.Abort()
		return
	}

	var data = map[string]string{
		"client_id":     os.Getenv("CLIENT_ID"),
		"client_secret": os.Getenv("CLIENT_SECRET"), // Add if your client requires a secret
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	keycloakResponse, err := utils.PostRequestUrlEncoded(tokenEndpoint, data, "")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	// Parse the response body
	var refreshTokenResponse utils.RefreshTokenResponse
	err = utils.UnmarshalResponse(keycloakResponse, &refreshTokenResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Return the new tokens
	c.JSON(http.StatusOK, refreshTokenResponse)
}

// Private functions

func createUserInKeyCloak(user *models.Users) (bool, error) {
	var keyCloakUser = utils.KeyCloakUserDto{
		Username:      user.Username,
		Enabled:       true,
		Email:         user.Email,
		EmailVerified: true,
		FirstName:     "John",
		LastName:      "Doe",
		Credentials: []utils.UserCredentialDto{
			{
				Type:      "password",
				Value:     user.Password,
				Temporary: false,
			},
		},
	}

	keyCloakToken, err := utils.GetKeyCloakAdminToken()

	if err != nil {
		return false, err
	}

	// Define the expected structure of the response
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var tokenResponse TokenResponse

	// Unmarshal the raw body into the ApiResponse struct
	err = utils.UnmarshalResponse(keyCloakToken, &tokenResponse)

	if err != nil {
		return false, errors.New("failed to get token")
	}
	//Prepare crate ser request
	url := fmt.Sprintf("%s/admin/realms/%s/users", os.Getenv("KEYCLOAK_HOST"), os.Getenv("REALM_NAME"))

	userData, err := json.Marshal(keyCloakUser)
	if err != nil {
		return false, err
	}

	_, _, err = utils.PostRequestJson(url, bytes.NewBuffer(userData), tokenResponse.AccessToken)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	return true, nil
}
