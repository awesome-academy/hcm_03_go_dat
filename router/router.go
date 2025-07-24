package router

import (
	"hotel-management/database"
	"hotel-management/internal/handler"
	"hotel-management/internal/handler/admin"
	"hotel-management/internal/middleware"
	"hotel-management/internal/repository"
	"hotel-management/internal/usecase"
	"hotel-management/internal/usecase/admin_usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	userRepository := repository.NewUserRepository(database.DB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	authUseCase := usecase.NewAuthUseCase(userRepository)
	authHandler := handler.NewAuthHandler(userUseCase, authUseCase)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh-token", authHandler.RefreshToken)
		authGroup.GET("/google/login", authHandler.GoogleLoginHandler)
		authGroup.GET("/google/callback", authHandler.GoogleCallbackHandler)
	}

	//Mail routes
	mailUseCase := usecase.NewMailUseCase(userRepository)
	mailHandler := handler.NewMailHandler(mailUseCase)
	mailGroup := r.Group("/mail")
	{
		mailGroup.POST("/smtp-verify", mailHandler.SendVerificationEmail)
		mailGroup.GET("/verify-account", mailHandler.ActiveAccountHandler)
		mailGroup.GET("/reset-password", mailHandler.ResetPassword)
	}

	//Admin route
	adminAuthUseCase := admin_usecase.NewAuthUseCase(userRepository)
	adminHandler := admin.NewAdminHandler(adminAuthUseCase)

	roomRepository := repository.NewRoomRepository(database.DB)
	reviewRepository := repository.NewReviewRepository(database.DB)
	bookingRepository := repository.NewBookingRepository(database.DB)
	roomUseCase := admin_usecase.NewRoomUseCase(roomRepository, bookingRepository, reviewRepository)
	roomHandler := admin.NewRoomHandler(roomUseCase)
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/", middleware.RequireLogin(), middleware.RequireRoles("admin"), adminHandler.AdminDashboard)
		adminGroup.GET("/login", adminHandler.AdminLoginPage)
		adminGroup.POST("/login", adminHandler.HandleLogin)
		adminGroup.GET("/logout", adminHandler.HandleLogout)
		adminGroup.GET("/rooms", middleware.RequireRoles("admin", "staff"), roomHandler.RoomManagementPage)
		adminGroup.GET("/rooms/create", middleware.RequireRoles("admin", "staff"), roomHandler.CreateRoomPage)
		adminGroup.POST("/rooms/create", middleware.RequireRoles("admin", "staff"), roomHandler.CreateRoom)
		adminGroup.GET("/rooms/:id", middleware.RequireRoles("admin", "staff"), roomHandler.RoomDetailPage)
		adminGroup.GET("/rooms/edit/:id", middleware.RequireRoles("admin", "staff"), roomHandler.EditRoomPage)
		adminGroup.POST("/rooms/edit/:id", middleware.RequireRoles("admin", "staff"), roomHandler.UpdateRoom)
		adminGroup.POST("/rooms/delete/:id", middleware.RequireRoles("admin", "staff"), roomHandler.DeleteRoom)
	}
	//User routes
	userHandler := handler.NewUserHandler(userUseCase)
	r.PUT("/users/update-profile", middleware.RequireAuth(userRepository), userHandler.UpdateProfile)

}
