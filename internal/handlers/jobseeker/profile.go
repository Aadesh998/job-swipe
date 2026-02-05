package jobseeker

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileInput struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ResumeURL    string `json:"resume_url"`
	Skills       string `json:"skills"`
	Experience   string `json:"experience"`
	Education    string `json:"education"`
	Bio          string `json:"bio"`
	IsOpenToWork bool   `json:"is_open_to_work"`
}

func CreateOrUpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "job_seeker" {
		response.Error(c, http.StatusForbidden, "Only job seekers can have a profile", nil)
		return
	}

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var profile models.JobSeekerProfile
	result := database.DB.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		profile = models.JobSeekerProfile{
			UserID:       userID.(uint),
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			ResumeURL:    input.ResumeURL,
			Skills:       input.Skills,
			Experience:   input.Experience,
			Education:    input.Education,
			Bio:          input.Bio,
			IsOpenToWork: input.IsOpenToWork,
		}
		if err := database.DB.Create(&profile).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create profile", err.Error())
			return
		}
		response.Created(c, "Profile created successfully", profile)
	} else {
		profile.FirstName = input.FirstName
		profile.LastName = input.LastName
		profile.ResumeURL = input.ResumeURL
		profile.Skills = input.Skills
		profile.Experience = input.Experience
		profile.Education = input.Education
		profile.Bio = input.Bio
		profile.IsOpenToWork = input.IsOpenToWork

		if err := database.DB.Save(&profile).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
			return
		}
		response.Success(c, "Profile updated successfully", profile)
	}
}

func GetProfile(c *gin.Context) {
	userID := c.Param("user_id") // Can request own or others

	// If no param, assume current user
	if userID == "" || userID == "me" {
		uid, _ := c.Get("user_id")
		var profile models.JobSeekerProfile
		if err := database.DB.Preload("Internships").Where("user_id = ?", uid).First(&profile).Error; err != nil {
			response.Error(c, http.StatusNotFound, "Profile not found", nil)
			return
		}
		response.Success(c, "Profile retrieved successfully", profile)
		return
	}

	// Fetch another user's profile (e.g., job provider viewing job seeker)
	var profile models.JobSeekerProfile
	if err := database.DB.Preload("User").Preload("Internships").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Profile not found", nil)
		return
	}
	
	// Hide sensitive user data from User preload if needed, or rely on json tags
	response.Success(c, "Profile retrieved successfully", profile)
}

func GetAllProfiles(c *gin.Context) {
	// Only for job providers or admin?
	role, _ := c.Get("role")
	if role != "job_provider" && role != "admin" {
		response.Error(c, http.StatusForbidden, "Unauthorized", nil)
		return
	}

	var profiles []models.JobSeekerProfile
	if err := database.DB.Preload("User").Find(&profiles).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch profiles", err.Error())
		return
	}

	response.Success(c, "Profiles retrieved successfully", profiles)
}
