package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/controllers"
	"github.com/leirbagxis/FreddyBot/internal/api/handlers"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func RegisterRoutes(r *gin.Engine, c *container.AppContainer) {
	loginRouter := r.Group("/auth")
	loginRouter.GET("/verify-token", handlers.VerifyJWTHandler())
	loginRouter.POST("/generate-token", handlers.GenerateJWTHandler(c))

	api := r.Group("/api")
	captionController := controllers.NewCaptionController(c)
	ButtonsController := controllers.NewButtonsController(c)
	permissionsController := controllers.NewPermissionController(c)

	api.Use(auth.AuthMiddlewareJWT())
	{
		api.GET("/ping", handlers.PingHandler(c))
		api.GET("/channel/:channelId", handlers.GetChannelHandler(c))

		api.PUT("/channel/:channelId/caption", captionController.UpdateDefaultCaptionController)
		api.PUT("/channel/:channelId/caption/permissions", permissionsController.UpdateMessagePermissionController)
		api.PUT("/channel/:channelId/buttons/permissions", permissionsController.UpdateButtonsPermissionController)

		api.POST("/channel/:channelId/buttons", ButtonsController.CreateDefaultButtonController)
		api.DELETE("/channel/:channelId/buttons/:buttonId", ButtonsController.DeleteDefaultButtonController)
		api.PUT("/channel/:channelId/buttons/:buttonId", ButtonsController.UpdateDefaultButtonController)
		api.PUT("/channel/:channelId/buttons/layout", ButtonsController.UpdateLayoutDefaultButtons)

	}
}
