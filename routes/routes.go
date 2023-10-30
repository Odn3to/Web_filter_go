package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"web-filter/controllers"
)

func HandleRequests(){
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	authRoutes := r.Group("/")
	authRoutes.Use(controllers.TokenValidationMiddleware)
	{
		authRoutes.POST("/new", controllers.CriaWebFilter)
		authRoutes.GET("/apply", controllers.ApplyWebFilter)
		authRoutes.GET("/status", controllers.GetStatusSquid)
		authRoutes.GET("/search/:searchValue", controllers.PesquisaWebFilter)
		authRoutes.GET("/search", controllers.PesquisaWebFilter)
		authRoutes.PUT("/edit/:id", controllers.EditarWebFilter)
		authRoutes.DELETE("/delete/:id", controllers.DeleteWebFilter)
	}

    r.Run(":8080")
}