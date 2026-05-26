package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leirbagxis/FreddyBot/internal/api/auth"
	"github.com/leirbagxis/FreddyBot/internal/api/controllers"
	admincontroller "github.com/leirbagxis/FreddyBot/internal/api/controllers/adminController"
	"github.com/leirbagxis/FreddyBot/internal/api/handlers"
	"github.com/leirbagxis/FreddyBot/internal/api/middleware"
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func RegisterRoutes(r *gin.Engine, c *container.AppContainer) {
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api")

	// Controladores
	authController := controllers.NewAuthController(c)
	captionController := controllers.NewCaptionController(c)
	ButtonsController := controllers.NewButtonsController(c)
	permissionsController := controllers.NewPermissionController(c)
	customCaptionController := controllers.NewCustomCaptionController(c)
	userController := controllers.NewUserController(c)
	channelController := controllers.NewChannelController(c)
	getALlUsers := admincontroller.NewUsersAdminController(c)
	configController := admincontroller.NewConfigController(c)
	mediaController := admincontroller.NewMediaController(c)
	auditController := admincontroller.NewAuditController(c)
	channelEventsController := admincontroller.NewChannelEventsController(c)

	// --- Rota de Login Unificada ---
	api.POST("/login", authController.Login)

	// --- Rotas Protegidas ---
	api.Use(auth.AuthMiddlewareJWT(c))
	{
		api.GET("/ping", handlers.PingHandler(c))
		api.GET("/me/channels", userController.GetUserChannelsController)
		api.GET("/user/info/:userParams", userController.GetUserInfo)
		api.POST("/channel/transfer", userController.TransferChannelController)

		// Rotas específicas de Canal (Com verificação de autorização)
		channelRoutes := api.Group("/channel/:channelId")
		channelRoutes.Use(auth.AuthorizeChannel(c))
		{
			channelRoutes.GET("", channelController.GetChannelByIDController)
			channelRoutes.DELETE("", channelController.DisconectChannel)
			channelRoutes.PUT("/caption", captionController.UpdateDefaultCaptionController)
			channelRoutes.PUT("/newpackcaption", captionController.UpdateNewPackCaptionController)
			channelRoutes.PUT("/reactions", captionController.UpdateReactionsController)
			channelRoutes.PUT("/reactions/active", permissionsController.UpdateReactionsActiveController)
			channelRoutes.PUT("/reactions/position", captionController.UpdateReactionPositionController)
			channelRoutes.PUT("/dynamic-links", permissionsController.UpdateDynamicLinksController)
			channelRoutes.PUT("/caption/permissions", permissionsController.UpdateMessagePermissionController)
			channelRoutes.PUT("/buttons/permissions", permissionsController.UpdateButtonsPermissionController)

			channelRoutes.POST("/buttons", ButtonsController.CreateDefaultButtonController)
			channelRoutes.DELETE("/buttons/:buttonId", ButtonsController.DeleteDefaultButtonController)
			channelRoutes.PUT("/buttons/:buttonId", ButtonsController.UpdateDefaultButtonController)
			channelRoutes.PUT("/buttons/layout", ButtonsController.UpdateLayoutDefaultButtons)

			channelRoutes.POST("/custom-captions", customCaptionController.CreateCustomCaptionController)
			channelRoutes.POST("/custom-captions/:captionId/buttons", customCaptionController.CreateCustomCaptionButtonController)
			channelRoutes.PUT("/custom-captions/:captionId", customCaptionController.UpdateCustomCaptionController)
			channelRoutes.PUT("/custom-captions/:captionId/layout", customCaptionController.UpdateCustomCaptionLayoutController)
			channelRoutes.PUT("/custom-captions/:captionId/buttons/:buttonId", customCaptionController.UpdateCustomCaptionButtonController)
			channelRoutes.DELETE("/custom-captions/:captionId", customCaptionController.DeleteCustomCaptionController)
			channelRoutes.DELETE("/custom-captions/:captionId/buttons/:buttonId", customCaptionController.DeleteCustomCaptionButtonController)

			channelRoutes.GET("/separator/:separatorId", channelController.GetSeparator)
		}
	}

	// --- Rotas Administrativas (Apenas Admin/Owner) ---
	adminRoute := api.Group("/admin")
	adminRoute.Use(auth.RequireRole(auth.RoleAdmin, auth.RoleOwner))
	{
		adminRoute.GET("/overview", getALlUsers.GetAdminOverview)
		adminRoute.GET("/users", getALlUsers.GetAllUsersAdminController)
		adminRoute.GET("/channels", channelController.GetAllChannelsController)
		adminRoute.POST("/notice", getALlUsers.SendNoticeAdminController)

		adminRoute.GET("/config", configController.GetConfig)
		adminRoute.PUT("/config", configController.UpdateConfig)

		adminRoute.GET("/media-proxy/:fileId", mediaController.GetMediaPreview)
		adminRoute.GET("/audit/checkbot", auditController.GetCheckBotAudit)
		adminRoute.GET("/logs", channelEventsController.List)
		adminRoute.POST("/audit/bulk-delete", auditController.BulkDeleteUserChannels)

		adminRoute.POST("/users/:userId/admin", getALlUsers.UpdateUserAdminController)
		adminRoute.POST("/users/:userId/blacklist", getALlUsers.UpdateUserBlacklistController)
	}
}
