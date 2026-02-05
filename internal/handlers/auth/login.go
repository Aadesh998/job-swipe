package auth

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"aron_project/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid email or password", nil)
		return
	}

	if !user.IsVerified {
		response.Error(c, http.StatusUnauthorized, "Please verify your email first", nil)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		response.Error(c, http.StatusBadRequest, "Invalid email or password", nil)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	response.Success(c, "Login successful", gin.H{"token": token})
}
