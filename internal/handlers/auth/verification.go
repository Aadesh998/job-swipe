package auth

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"aron_project/internal/utils"
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

	jwtToken, _ := utils.GenerateToken(user.ID, user.Email, user.Role)

	response.Success(c, "Email verified successfully", gin.H{"token": jwtToken})
}
