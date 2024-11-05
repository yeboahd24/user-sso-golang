package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yeboahd24/user-sso/handler"
)

//
// func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
// 	auth := r.Group("/auth")
// 	{
// 		auth.POST("/register", authHandler.Register)
// 		auth.GET("/google/login", authHandler.GoogleLogin)
// 		auth.GET("/google/callback", authHandler.GoogleCallback)
// 	}
// }

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, tmplHandler *handler.TemplateHandler) {
	// Load HTML templates
	r.LoadHTMLGlob("template/*.html")

	// Public routes
	r.GET("/", tmplHandler.LoginPage)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/google/login", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.GET("/verify", authHandler.VerifySession)
		auth.POST("/logout", authHandler.Logout)
	}
}
