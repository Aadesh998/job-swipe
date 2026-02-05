package job

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetJobApplicants retrieves all applicants for a specific job
func GetJobApplicants(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("job_id")

	// Verify job ownership
	var job models.Job
	if err := database.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Job not found", nil)
		return
	}

	if job.Company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this job posting", nil)
		return
	}

	var applications []models.Application
	// Preload User, JobSeekerProfile and its related data (Internships)
	if err := database.DB.Preload("User").Preload("User.JobSeekerProfile").Preload("User.JobSeekerProfile.Internships").Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch applicants", err.Error())
		return
	}

	response.Success(c, "Applicants retrieved successfully", applications)
}

type UpdateApplicationStatusInput struct {
	Status string `json:"status" binding:"required,oneof=applied reviewing interviewed rejected hired"`
}

// UpdateApplicationStatus allows job provider to change status (e.g., Shortlist, Reject)
func UpdateApplicationStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	applicationID := c.Param("application_id")

	var application models.Application
	if err := database.DB.Preload("Job").Preload("Job.Company").First(&application, applicationID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Application not found", nil)
		return
	}

	// Verify ownership via Job -> Company -> UserID
	if application.Job.Company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this job application", nil)
		return
	}

	var input UpdateApplicationStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	application.Status = input.Status
	if err := database.DB.Save(&application).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update status", err.Error())
		return
	}

	response.Success(c, "Application status updated", application)
}
