package tiny_url

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	tinyUrlHandler := NewHandler()
	router.POST("/shorten", tinyUrlHandler.CreateTinyUrl)
	router.GET("/:shortcode", tinyUrlHandler.SingleTinyUrl)
	router.GET("/:shortcode/stats", tinyUrlHandler.StatSingleTinyUrl)
}
