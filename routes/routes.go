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

	r.POST("/new", controllers.CriaWebFilter)
	r.GET("/apply", controllers.ApplyWebFilter)
	r.GET("/search/:searchValue", controllers.PesquisaWebFilter)
	r.GET("/search", controllers.PesquisaWebFilter)
	r.PUT("/edit/:id", controllers.EditarWebFilter)
	r.DELETE("/delete/:id", controllers.DeleteWebFilter)

    r.Run(":8080")
}