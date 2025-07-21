package routing

import (
	"xyz-multifinance/config"
	"xyz-multifinance/internal/delivery/http"
	"xyz-multifinance/internal/repository"
	"xyz-multifinance/internal/usecase"
	"xyz-multifinance/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	// Middleware global
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Inject dependencies
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	limitRepo := repository.NewLimitRepository(db)
	customerRepo := repository.NewCustomerRepository(db)

	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewAuthHandler(userUC)

	customerUC := usecase.NewCustomerUsecase(customerRepo, userRepo)
	customerHandler := http.NewCustomerHandler(customerUC)

	limitUC := usecase.NewLimitUsecase(limitRepo, transactionRepo)
	limitHandler := http.NewLimitHandler(limitUC)

	transactionUC := usecase.NewTransactionUsecase(transactionRepo, limitRepo, customerRepo, db)
	transactionHandler := http.NewTransactionHandler(transactionUC)

	// Public routes
	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api.POST("/login", userHandler.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.Auth(cfg))

	protected.POST("/users", userHandler.CreateUser)

	// Customer routes
	protected.POST("/customers", middleware.AdminOnly(), customerHandler.CreateCustomer)
	protected.GET("/customers/:nik", customerHandler.GetCustomerByNIK)
	protected.PUT("/customers/:nik", customerHandler.UpdateCustomer)
	protected.DELETE("/customers/:nik", customerHandler.DeleteCustomer)

	// Limit routes
	protected.POST("/limits", limitHandler.CreateLimit)
	protected.PUT("/limits/:id", limitHandler.UpdateLimit)
	protected.DELETE("/limits/:id", limitHandler.DeleteLimit)
	protected.GET("/limits/:id", limitHandler.GetLimitByID)
	protected.GET("/limits/customer/:customer_id", limitHandler.GetLimitsByCustomerID)
	protected.GET("/limits/customer/:customer_id/tenor/:tenor", limitHandler.GetLimitByCustomerAndTenor)

	// Transaction routes
	protected.POST("/transactions", transactionHandler.CreateTransaction)
	protected.GET("/transactions/:nik", transactionHandler.GetTransactionsByCustomer)

	// Handle no route/method
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"error": "Method not allowed"})
	})
}
