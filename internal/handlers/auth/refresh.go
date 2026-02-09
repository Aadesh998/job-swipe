package auth

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"job_swipe/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	claims, err := utils.ValidateToken(input.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token", err.Error())
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Invalid token claims", nil)
		return
	}
	userID := uint(userIDFloat)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, "User not found", nil)
		return
	}

	newTokenPair, err := utils.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate tokens", err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", newTokenPair)
}
