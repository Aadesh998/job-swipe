package job

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JobInput struct {
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description" binding:"required"`
	Requirements      string `json:"requirements"`
	YearsOfExperience string `json:"years_of_experience"`
	Field             string `json:"field" binding:"required"`
	Location          string `json:"location"`
	Type              string `json:"type" binding:"required"`
	Stipend           string `json:"stipend"`
	SalaryRange       string `json:"salary_range"`
	Status            string `json:"status"`
}

func CreateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	companyID := c.Param("company_id")

	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	if company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this company", nil)
		return
	}

	var input JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	job := models.Job{
		CompanyID:         company.ID,
		Title:             input.Title,
		Description:       input.Description,
		Requirements:      input.Requirements,
		YearsOfExperience: input.YearsOfExperience,
		Field:             input.Field,
		Location:          input.Location,
		Type:              input.Type,
		Stipend:           input.Stipend,
		SalaryRange:       input.SalaryRange,
		Status:            "open",
	}

	if err := database.DB.Create(&job).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create job", err.Error())
		return
	}

	response.Created(c, "Job created successfully", job)
}

func UpdateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("job_id")

	var job models.Job
	if err := database.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Job not found", nil)
		return
	}

	if job.Company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this job posting", nil)
		return
	}

	var input JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	job.Title = input.Title
	job.Description = input.Description
	job.Requirements = input.Requirements
	job.YearsOfExperience = input.YearsOfExperience
	job.Field = input.Field
	job.Location = input.Location
	job.Type = input.Type
	job.Stipend = input.Stipend
	job.SalaryRange = input.SalaryRange
	if input.Status != "" {
		job.Status = input.Status
	}

	if err := database.DB.Save(&job).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update job", err.Error())
		return
	}

	response.Success(c, "Job updated successfully", job)
}

func DeleteJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("job_id")

	var job models.Job
	if err := database.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Job not found", nil)
		return
	}

	if job.Company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this job posting", nil)
		return
	}

	if err := database.DB.Delete(&job).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete job", err.Error())
		return
	}

	response.Success(c, "Job deleted successfully", nil)
}

func GetJob(c *gin.Context) {
	jobID := c.Param("job_id")

	var job models.Job
	if err := database.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Job not found", nil)
		return
	}

	response.Success(c, "Job details", job)
}

func GetCompanyJobs(c *gin.Context) {
	companyID := c.Param("company_id")

	var jobs []models.Job
	if err := database.DB.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch jobs", err.Error())
		return
	}

	response.Success(c, "Company jobs", jobs)
}

func GetAllJobs(c *gin.Context) {
	var jobs []models.Job
	// Could add filtering here later
	if err := database.DB.Preload("Company").Where("status = ?", "open").Find(&jobs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch jobs", err.Error())
		return
	}
	response.Success(c, "All open jobs", jobs)
}
