package app

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/auth", func(c *gin.Context) {
	})

	r.Use(func(context *gin.Context) {
		auth := context.GetHeader("auth")
		if auth == "" {

		}
	})
	r.GET("/list", func(c *gin.Context) {
	})
	r.GET("/exist/:uid", func(c *gin.Context) {
	})
	r.GET("/count", func(c *gin.Context) {
	})
	r.GET("/info/:uid", func(c *gin.Context) {
	})
	r.POST("/del/:uid", func(c *gin.Context) {
	})
	r.POST("/online/:uid", func(c *gin.Context) {
	})
	if err := r.Run(":8080");err!=nil{

	}
}