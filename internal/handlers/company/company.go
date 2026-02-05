package company

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompanyInput struct {
	CompanyName string `json:"company_name" binding:"required"`
	CompanySize string `json:"company_size" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Industry    string `json:"industry" binding:"required"`
	Services    string `json:"services"`
	LogoURL     string `json:"logo_url"`
}

// CreateCompany creates a new company profile for the user
func CreateCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	role, _ := c.Get("role")
	if role != "job_provider" {
		response.Error(c, http.StatusForbidden, "Only job providers can create a company profile", nil)
		return
	}

	var input CompanyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	company := models.Company{
		UserID:            userID.(uint),
		CompanyName:       input.CompanyName,
		CompanySize:       input.CompanySize,
		Location:          input.Location,
		Website:           input.Website,
		Description:       input.Description,
		Industry:          input.Industry,
		Services:          input.Services,
		LogoURL:           input.LogoURL,
		IsProfileComplete: true,
	}

	if err := database.DB.Create(&company).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create company profile", err.Error())
		return
	}
	response.Created(c, "Company profile created successfully", company)
}

// GetUserCompanies returns all companies owned by the user
func GetUserCompanies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var companies []models.Company
	if err := database.DB.Where("user_id = ?", userID).Preload("Products").Find(&companies).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch companies", err.Error())
		return
	}

	response.Success(c, "Companies retrieved successfully", companies)
}

// GetCompany returns a specific company details
func GetCompany(c *gin.Context) {
	id := c.Param("id")
	var company models.Company
	if err := database.DB.Preload("Products").First(&company, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	response.Success(c, "Company details", company)
}

// UpdateCompany updates a specific company
func UpdateCompany(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	var company models.Company
	
	if err := database.DB.First(&company, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	// Verify ownership
	if company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this company", nil)
		return
	}

	var input CompanyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	company.CompanyName = input.CompanyName
	company.CompanySize = input.CompanySize
	company.Location = input.Location
	company.Website = input.Website
	company.Description = input.Description
	company.Industry = input.Industry
	company.Services = input.Services
	company.LogoURL = input.LogoURL

	if err := database.DB.Save(&company).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update company", err.Error())
		return
	}

	response.Success(c, "Company updated successfully", company)
}