package provider

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileInput struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Title         string `json:"title"`
	ContactNumber string `json:"contact_number"`
	Bio           string `json:"bio"`
}

func CreateOrUpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "job_provider" {
		response.Error(c, http.StatusForbidden, "Only job providers can have a provider profile", nil)
		return
	}

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var profile models.JobProviderProfile
	result := database.DB.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		profile = models.JobProviderProfile{
			UserID:        userID.(uint),
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			Title:         input.Title,
			ContactNumber: input.ContactNumber,
			Bio:           input.Bio,
		}
		if err := database.DB.Create(&profile).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create profile", err.Error())
			return
		}
		response.Created(c, "Profile created successfully", profile)
	} else {
		profile.FirstName = input.FirstName
		profile.LastName = input.LastName
		profile.Title = input.Title
		profile.ContactNumber = input.ContactNumber
		profile.Bio = input.Bio

		if err := database.DB.Save(&profile).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
			return
		}
		response.Success(c, "Profile updated successfully", profile)
	}
}

func GetProfile(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" || userID == "me" {
		uid, _ := c.Get("user_id")
		var profile models.JobProviderProfile
		if err := database.DB.Where("user_id = ?", uid).First(&profile).Error; err != nil {
			response.Error(c, http.StatusNotFound, "Profile not found", nil)
			return
		}
		response.Success(c, "Profile retrieved successfully", profile)
		return
	}

	var profile models.JobProviderProfile
	if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Profile not found", nil)
		return
	}
	response.Success(c, "Profile retrieved successfully", profile)
}
