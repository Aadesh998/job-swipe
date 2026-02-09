package jobseeker

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"job_swipe/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InternshipInput struct {
	Company      string    `json:"company" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsCurrent    bool      `json:"is_current"`
	Location     string    `json:"location"`
	Technologies string    `json:"technologies"`
}

func AddInternship(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var profile models.JobSeekerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Please create a profile first", nil)
		return
	}

	var input InternshipInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	internship := models.Internship{
		JobSeekerProfileID: profile.ID,
		Company:            input.Company,
		Role:               input.Role,
		Description:        input.Description,
		StartDate:          input.StartDate,
		EndDate:            input.EndDate,
		IsCurrent:          input.IsCurrent,
		Location:           input.Location,
		Technologies:       input.Technologies,
	}

	if err := database.DB.Create(&internship).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to add internship", err.Error())
		return
	}

	response.Created(c, "Internship added successfully", internship)
}

func UpdateInternship(c *gin.Context) {
	userID, _ := c.Get("user_id")
	internshipID := c.Param("internship_id")

	var internship models.Internship
	if err := database.DB.First(&internship, internshipID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Internship not found", nil)
		return
	}

	var profile models.JobSeekerProfile
	if err := database.DB.First(&profile, internship.JobSeekerProfileID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Profile not found", nil)
		return
	}

	if profile.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this record", nil)
		return
	}

	var input InternshipInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	internship.Company = input.Company
	internship.Role = input.Role
	internship.Description = input.Description
	internship.StartDate = input.StartDate
	internship.EndDate = input.EndDate
	internship.IsCurrent = input.IsCurrent
	internship.Location = input.Location
	internship.Technologies = input.Technologies

	if err := database.DB.Save(&internship).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update internship", err.Error())
		return
	}

	response.Success(c, "Internship updated successfully", internship)
}

func DeleteInternship(c *gin.Context) {
	userID, _ := c.Get("user_id")
	internshipID := c.Param("internship_id")

	var internship models.Internship
	if err := database.DB.First(&internship, internshipID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Internship not found", nil)
		return
	}

	var profile models.JobSeekerProfile
	if err := database.DB.First(&profile, internship.JobSeekerProfileID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Profile not found", nil)
		return
	}

	if profile.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this record", nil)
		return
	}

	if err := database.DB.Delete(&internship).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete internship", err.Error())
		return
	}

	response.Success(c, "Internship deleted successfully", nil)
}
