package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/controllers"
	admincontroller "github.com/leirbagxis/FreddyBot/internal/api/controllers/adminController"
	webappauthcontroller "github.com/leirbagxis/FreddyBot/internal/api/controllers/webAppAuthController"
	"github.com/leirbagxis/FreddyBot/internal/api/handlers"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func RegisterRoutes(r *gin.Engine, c *container.AppContainer) {
	// loginRouter := r.Group("/auth")
	// loginRouter.GET("/verify-token", handlers.VerifyJWTHandler())
	// loginRouter.POST("/generate-token", handlers.GenerateJWTHandler(c))

	api := r.Group("/api")
	captionController := controllers.NewCaptionController(c)
	ButtonsController := controllers.NewButtonsController(c)
	permissionsController := controllers.NewPermissionController(c)
	customCaptionController := controllers.NewCustomCaptionController(c)
	separatorController := controllers.NewSeparatorController(c)
	webAppAuthController := webappauthcontroller.NewWebAppAuthController(c)
	userController := controllers.NewUserController(c)
	channelController := controllers.NewChannelController(c)

	api.POST("/auth", webAppAuthController.ReceiveAuthController)
	api.POST("/me/channels", webAppAuthController.ReceiveAuthMeChannelsController)
	api.POST("/admin/dash", webAppAuthController.AdminAuthController)

	getALlUsers := admincontroller.NewUsersAdminController(c)

	api.Use(auth.AuthMiddlewareJWT(c))
	{
		api.GET("/ping", handlers.PingHandler(c))
		api.GET("/channel/:channelId", handlers.GetChannelHandler(c))

		api.PUT("/channel/:channelId/caption", captionController.UpdateDefaultCaptionController)
		api.PUT("/channel/:channelId/newpackcaption", captionController.UpdateNewPackCaptionController)
		api.PUT("/channel/:channelId/caption/permissions", permissionsController.UpdateMessagePermissionController)
		api.PUT("/channel/:channelId/buttons/permissions", permissionsController.UpdateButtonsPermissionController)

		api.POST("/channel/:channelId/buttons", ButtonsController.CreateDefaultButtonController)
		api.DELETE("/channel/:channelId/buttons/:buttonId", ButtonsController.DeleteDefaultButtonController)
		api.PUT("/channel/:channelId/buttons/:buttonId", ButtonsController.UpdateDefaultButtonController)
		api.PUT("/channel/:channelId/buttons/layout", ButtonsController.UpdateLayoutDefaultButtons)

		api.POST("/channel/:channelId/custom-captions", customCaptionController.CreateCustomCaptionController)
		api.POST("/channel/:channelId/custom-captions/:captionId/buttons", customCaptionController.CreateCustomCaptionButtonController)
		api.PUT("/channel/:channelId/custom-captions/:captionId", customCaptionController.UpdateCustomCaptionController)
		api.PUT("/channel/:channelId/custom-captions/:captionId/layout", customCaptionController.UpdateCustomCaptionLayoutController)
		api.PUT("/channel/:channelId/custom-captions/:captionId/buttons/:buttonId", customCaptionController.UpdateCustomCaptionButtonController)
		api.DELETE("/channel/:channelId/custom-captions/:captionId", customCaptionController.DeleteCustomCaptionController)
		api.DELETE("/channel/:channelId/custom-captions/:captionId/buttons/:buttonId", customCaptionController.DeleteCustomCaptionButtonController)

		api.GET("/channel/:channelId/separator/:separatorId", separatorController.GetSeparator)
		api.DELETE("/channel/disconect", channelController.DisconectChannel)

		api.GET("/user/info/:userParams", userController.GetUserInfo)
		api.POST("/channel/transfer", userController.TransferChannelController)
		api.POST("/admin/notice", getALlUsers.SendNoticeAdminController)
	}

	adminRoute := r.Group("/admin/api")

	adminRoute.Use(auth.AuthMiddlewareJWT(c))
	{
		adminRoute.GET("/users", getALlUsers.GetAllUsersAdminController)
	}
}
