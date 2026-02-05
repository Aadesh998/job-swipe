package auth

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"aron_project/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConfig *oauth2.Config

func getGoogleConfig() *oauth2.Config {
	if googleConfig == nil {
		googleConfig = &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}
	return googleConfig
}

func GoogleLogin(c *gin.Context) {
	url := getGoogleConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "Code not found", nil)
		return
	}

	token, err := getGoogleConfig().Exchange(context.Background(), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to exchange token", err.Error())
		return
	}

	// Fetch user info using the token
	client := getGoogleConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get user info", err.Error())
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to decode user info", err.Error())
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", googleUser.Email).First(&user)

	if result.Error != nil {
		// User not found, create new user
		user = models.User{
			Email:      googleUser.Email,
			GoogleID:   googleUser.ID,
			IsVerified: true, // Google emails are verified
			Role:       "job_seeker",
			AvatarURL:  googleUser.Picture,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}
	} else {
		// User found, update GoogleID if missing
		if user.GoogleID == "" {
			user.GoogleID = googleUser.ID
			user.AvatarURL = googleUser.Picture
			user.IsVerified = true // Trust Google
			database.DB.Save(&user)
		}
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	response.Success(c, "Login successful", gin.H{
		"token": jwtToken,
		"user":  user,
	})
}
