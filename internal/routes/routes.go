package routes

import (
	"aron_project/internal/handlers"
	"aron_project/internal/handlers/auth"
	"aron_project/internal/handlers/chat"
	"aron_project/internal/handlers/company"
	"aron_project/internal/handlers/job"
	"aron_project/internal/handlers/jobseeker"
	"aron_project/internal/handlers/provider"
	"aron_project/internal/middleware"
	"aron_project/internal/response"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	go chat.GlobalHub.Run()

	r.Use(middleware.RateLimitMiddleware())

	r.GET("/health", handlers.HealthCheck)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", auth.Signup)
		authGroup.GET("/verify-email", auth.VerifyEmail)
		authGroup.POST("/login", auth.Login)
		authGroup.GET("/google/login", auth.GoogleLogin)
		authGroup.GET("/google/callback", auth.GoogleCallback)
		authGroup.POST("/forgot-password", auth.ForgotPassword)
		authGroup.POST("/reset-password", auth.ResetPassword)
		authGroup.GET("/generate-password", auth.GetPasswordGenerator)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			role, _ := c.Get("role")
			response.Success(c, "Profile data", gin.H{
				"user_id": userID,
				"email":   email,
				"role":    role,
			})
		})

		// Job Provider Personal Profile Routes
		jobProviderGroup := api.Group("/job-provider")
		{
			jobProviderGroup.POST("/profile", provider.CreateOrUpdateProfile)
			jobProviderGroup.GET("/profile", provider.GetProfile)
			jobProviderGroup.GET("/profile/:user_id", provider.GetProfile)
		}

		// Company routes
		companyGroup := api.Group("/companies")
		{
			companyGroup.POST("/", company.CreateCompany)
			companyGroup.GET("/", company.GetUserCompanies)
			companyGroup.GET("/:id", company.GetCompany)
			companyGroup.PUT("/:id", company.UpdateCompany)

			// Product routes nested under company
			companyGroup.POST("/:id/products", company.AddProduct)
			companyGroup.GET("/:id/products", company.GetCompanyProducts)

			// Job routes nested under company
			companyGroup.POST("/:id/jobs", job.CreateJob)
			companyGroup.GET("/:id/jobs", job.GetCompanyJobs)
		}

		// Direct product manipulation
		productGroup := api.Group("/products")
		{
			productGroup.GET("/:product_id", company.GetProduct)
			productGroup.PUT("/:product_id", company.UpdateProduct)
			productGroup.DELETE("/:product_id", company.DeleteProduct)
		}

		// Job routes (general)
		jobGroup := api.Group("/jobs")
		{
			jobGroup.GET("/", job.GetAllJobs)
			jobGroup.GET("/:job_id", job.GetJob)
			jobGroup.PUT("/:job_id", job.UpdateJob)
			jobGroup.DELETE("/:job_id", job.DeleteJob)

			// Applicant routes
			jobGroup.GET("/:job_id/applicants", job.GetJobApplicants)
			jobGroup.PUT("/applications/:application_id/status", job.UpdateApplicationStatus)
		}

		// Job Seeker routes
		jobSeekerGroup := api.Group("/job-seeker")
		{
			jobSeekerGroup.POST("/profile", jobseeker.CreateOrUpdateProfile)
			jobSeekerGroup.GET("/profile", jobseeker.GetProfile)
			jobSeekerGroup.GET("/profile/:user_id", jobseeker.GetProfile)
			jobSeekerGroup.GET("/profiles", jobseeker.GetAllProfiles)

			// Internship routes
			jobSeekerGroup.POST("/internships", jobseeker.AddInternship)
			jobSeekerGroup.PUT("/internships/:internship_id", jobseeker.UpdateInternship)
			jobSeekerGroup.DELETE("/internships/:internship_id", jobseeker.DeleteInternship)

			// Job Discovery (Tinder-like)
			jobSeekerGroup.GET("/jobs/discovery", jobseeker.GetJobsForSwipe)
			jobSeekerGroup.POST("/jobs/swipe", jobseeker.SwipeJob)
			jobSeekerGroup.GET("/jobs/search", jobseeker.SearchJobs)
		}

		// Chat routes
		chatGroup := api.Group("/chat")
		{
			chatGroup.GET("/ws", func(c *gin.Context) {
				chat.ServeWs(chat.GlobalHub, c)
			})
			chatGroup.POST("/send", chat.SendMessage)
			chatGroup.GET("/history/:user_id", chat.GetChatHistory)
			chatGroup.GET("/conversations", chat.GetConversations)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				response.Success(c, "Welcome Admin", nil)
			})
		}
	}
}
