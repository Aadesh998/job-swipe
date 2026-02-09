package auth

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"job_swipe/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, "Token required", nil)
		return
	}

	var user models.User
	if err := database.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid token", nil)
		return
	}

	user.IsVerified = true
	user.VerificationToken = ""
	database.DB.Save(&user)

	jwtToken, _ := utils.GenerateTokenPair(user.ID, user.Email, user.Role)

	response.Success(c, "Email verified successfully", gin.H{"token": jwtToken})
}
