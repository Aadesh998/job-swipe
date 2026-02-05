package auth

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"aron_project/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response.Success(c, "If this email is registered, you will receive a reset link", nil)
		return
	}

	resetToken, _ := utils.GenerateRandomString(32)
	expiry := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = resetToken
	user.ResetTokenExpiry = &expiry
	database.DB.Save(&user)

	resetLink := "http://localhost:5900/auth/reset-password-page?token=" + resetToken
	emailBody := "Click here to reset your password: <a href='" + resetLink + "'>Reset Password</a>"
	go utils.SendEmail(user.Email, "Reset Password", emailBody)

	response.Success(c, "If this email is registered, you will receive a reset link", nil)
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("password_reset_token = ?", input.Token).First(&user).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid or expired token", nil)
		return
	}

	if user.ResetTokenExpiry != nil && time.Now().After(*user.ResetTokenExpiry) {
		response.Error(c, http.StatusBadRequest, "Token expired", nil)
		return
	}

	hashedPassword, _ := utils.HashPassword(input.NewPassword)
	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.ResetTokenExpiry = nil
	database.DB.Save(&user)

	response.Success(c, "Password reset successfully", nil)
}

func GetPasswordGenerator(c *gin.Context) {
	password, _ := utils.GenerateRandomString(12)
	response.Success(c, "Password generated", gin.H{"generated_password": password})
}
