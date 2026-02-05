package company

import (
	"aron_project/internal/database"
	"aron_project/internal/models"
	"aron_project/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

// AddProduct adds a product to a specific company
func AddProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	companyID := c.Param("id")

	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	if company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this company", nil)
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	product := models.Product{
		CompanyID:   company.ID,
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to add product", err.Error())
		return
	}

	response.Created(c, "Product added successfully", product)
}

// GetCompanyProducts retrieves all products for a specific company
func GetCompanyProducts(c *gin.Context) {
	companyID := c.Param("id")
	
	// Check if company exists
	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	var products []models.Product
	if err := database.DB.Where("company_id = ?", companyID).Find(&products).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch products", err.Error())
		return
	}

	response.Success(c, "Products retrieved successfully", products)
}

// UpdateProduct updates a product's details
func UpdateProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	productID := c.Param("product_id")

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", nil)
		return
	}

	// Verify ownership via company
	var company models.Company
	if err := database.DB.First(&company, product.CompanyID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Associated company not found", nil)
		return
	}

	if company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this product", nil)
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.ImageURL = input.ImageURL
	product.Price = input.Price

	if err := database.DB.Save(&product).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	response.Success(c, "Product updated successfully", product)
}

// DeleteProduct removes a product
func DeleteProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	productID := c.Param("product_id")

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", nil)
		return
	}

	// Verify ownership
	var company models.Company
	if err := database.DB.First(&company, product.CompanyID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Associated company not found", nil)
		return
	}

	if company.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "You do not own this product", nil)
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}

// GetProduct retrieves a single product
func GetProduct(c *gin.Context) {
	id := c.Param("product_id")
	
	// Validate ID is numeric
	if _, err := strconv.Atoi(id); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", nil)
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}
