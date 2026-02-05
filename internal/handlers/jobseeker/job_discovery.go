package jobseeker

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetJobsForSwipe returns jobs that the user hasn't swiped on yet
func GetJobsForSwipe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Get IDs of jobs already swiped
	var swipedJobIDs []uint

	database.DB.Model(&models.JobSwipe{}).Where("user_id = ?", userID).Pluck("job_id", &swipedJobIDs)

	query := database.DB.Preload("Company").Where("status = ?", "open")

	if len(swipedJobIDs) > 0 {
		query = query.Where("id NOT IN ?", swipedJobIDs)
	}

	// Optional: Filter by preferences if stored in profile
	// For now, just return random or latest open jobs
	var jobs []models.Job
	if err := query.Limit(10).Order("created_at desc").Find(&jobs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch jobs", err.Error())
		return
	}

	response.Success(c, "Jobs for you", jobs)
}

type SwipeInput struct {
	JobID  uint   `json:"job_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=like pass"`
}

// SwipeJob handles the user's action on a job (like/apply or pass)
func SwipeJob(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input SwipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	// Check if already swiped
	var existingSwipe models.JobSwipe
	if err := database.DB.Where("user_id = ? AND job_id = ?", userID, input.JobID).First(&existingSwipe).Error; err == nil {
		response.Error(c, http.StatusBadRequest, "You have already swiped on this job", nil)
		return
	}

	// Use transaction
	tx := database.DB.Begin()

	swipe := models.JobSwipe{
		UserID: userID.(uint),
		JobID:  input.JobID,
		Action: input.Action,
	}

	if err := tx.Create(&swipe).Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "Failed to save action", err.Error())
		return
	}

	message := "Job passed"
	if input.Action == "like" {
		application := models.Application{
			JobID:  input.JobID,
			UserID: userID.(uint),
			Status: "applied",
		}
		if err := tx.Create(&application).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "Failed to create application", err.Error())
			return
		}
		message = "Job applied/liked successfully"
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Error(c, http.StatusInternalServerError, "Transaction commit failed", err.Error())
		return
	}

	response.Success(c, message, nil)
}

func SearchJobs(c *gin.Context) {
	field := c.Query("field")
	location := c.Query("location")
	jobType := c.Query("type")

	query := database.DB.Preload("Company").Where("status = ?", "open")

	if field != "" {
		query = query.Where("field ILIKE ?", "%"+field+"%")
	}
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if jobType != "" {
		query = query.Where("type = ?", jobType)
	}

	var jobs []models.Job
	if err := query.Find(&jobs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to search jobs", err.Error())
		return
	}

	response.Success(c, "Search results", jobs)
}
