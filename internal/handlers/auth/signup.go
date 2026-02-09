package auth

import (
	"fmt"
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"job_swipe/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=job_seeker job_provider admin"`
}

func Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	if input.Role == "" {
		input.Role = "job_seeker"
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		response.Error(c, http.StatusBadRequest, "Email already registered", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	verificationToken, _ := utils.GenerateRandomString(32)

	user := models.User{
		Email:             input.Email,
		Password:          hashedPassword,
		VerificationToken: verificationToken,
		IsVerified:        false,
		Role:              input.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	verifyLink := fmt.Sprintf("%s/auth/verify-email?token=%s", c.Request.Host, verificationToken)
	emailBody := "Click here to verify your email: <a href='" + verifyLink + "'>Verify Email</a>"
	go utils.SendEmail(user.Email, "Verify your email", emailBody)

	response.Created(c, "User created. Please check your email to verify.", nil)
}
