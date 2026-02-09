package routes

import (
	"job_swipe/internal/handlers"
	"job_swipe/internal/handlers/auth"
	"job_swipe/internal/handlers/chat"
	"job_swipe/internal/handlers/company"
	"job_swipe/internal/handlers/job"
	"job_swipe/internal/handlers/jobseeker"
	"job_swipe/internal/handlers/provider"
	"job_swipe/internal/middleware"
	"job_swipe/internal/response"

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
		authGroup.POST("/refresh", auth.RefreshToken)
		authGroup.GET("/google/login", auth.GoogleLogin)
		authGroup.GET("/google/callback", auth.GoogleCallback)
		authGroup.POST("/forgot-password", auth.ForgotPassword)
		authGroup.POST("/reset-password", auth.ResetPassword)
		authGroup.GET("/generate-password", auth.GetPasswordGenerator)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	v1 := api.Group("/v1")
	{
		v1.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			role, _ := c.Get("role")
			response.Success(c, "Profile data", gin.H{
				"user_id": userID,
				"email":   email,
				"role":    role,
			})
		})

		jobProviderGroup := v1.Group("/job-provider")
		{
			jobProviderGroup.POST("/profile", provider.CreateOrUpdateProfile)
			jobProviderGroup.GET("/profile", provider.GetProfile)
			jobProviderGroup.GET("/profile/:user_id", provider.GetProfile)
		}

		companyGroup := v1.Group("/companies")
		{
			companyGroup.POST("/", company.CreateCompany)
			companyGroup.GET("/", company.GetUserCompanies)
			companyGroup.GET("/:id", company.GetCompany)
			companyGroup.PUT("/:id", company.UpdateCompany)
			companyGroup.POST("/:id/products", company.AddProduct)
			companyGroup.GET("/:id/products", company.GetCompanyProducts)
			companyGroup.POST("/:id/jobs", job.CreateJob)
			companyGroup.GET("/:id/jobs", job.GetCompanyJobs)
		}

		productGroup := v1.Group("/products")
		{
			productGroup.GET("/:product_id", company.GetProduct)
			productGroup.PUT("/:product_id", company.UpdateProduct)
			productGroup.DELETE("/:product_id", company.DeleteProduct)
		}

		jobGroup := v1.Group("/jobs")
		{
			jobGroup.GET("/", job.GetAllJobs)
			jobGroup.GET("/:job_id", job.GetJob)
			jobGroup.PUT("/:job_id", job.UpdateJob)
			jobGroup.DELETE("/:job_id", job.DeleteJob)
			jobGroup.GET("/:job_id/applicants", job.GetJobApplicants)
			jobGroup.PUT("/applications/:application_id/status", job.UpdateApplicationStatus)
		}

		jobSeekerGroup := v1.Group("/job-seeker")
		{
			jobSeekerGroup.POST("/profile", jobseeker.CreateOrUpdateProfile)
			jobSeekerGroup.GET("/profile", jobseeker.GetProfile)
			jobSeekerGroup.GET("/profile/:user_id", jobseeker.GetProfile)
			jobSeekerGroup.GET("/profiles", jobseeker.GetAllProfiles)
			jobSeekerGroup.POST("/internships", jobseeker.AddInternship)
			jobSeekerGroup.PUT("/internships/:internship_id", jobseeker.UpdateInternship)
			jobSeekerGroup.DELETE("/internships/:internship_id", jobseeker.DeleteInternship)
			jobSeekerGroup.GET("/jobs/discovery", jobseeker.GetJobsForSwipe)
			jobSeekerGroup.POST("/jobs/swipe", jobseeker.SwipeJob)
			jobSeekerGroup.GET("/jobs/search", jobseeker.SearchJobs)
		}

		chatGroup := v1.Group("/chat")
		{
			chatGroup.GET("/ws", func(c *gin.Context) {
				chat.ServeWs(chat.GlobalHub, c)
			})
			chatGroup.POST("/send", chat.SendMessage)
			chatGroup.GET("/history/:user_id", chat.GetChatHistory)
			chatGroup.GET("/conversations", chat.GetConversations)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				response.Success(c, "Welcome Admin", nil)
			})
		}
	}
}
