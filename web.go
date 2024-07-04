package main

// gin
import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/v1/chat/completions", OpenAIRequest)

	return r
}

func OpenAIRequest(c *gin.Context) {
	var chatCompletion OpenAIChatCompletion

	token := c.GetHeader("Authorization")
	if token != "Bearer "+ConfigInstance.AuthKey {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	err := c.ShouldBind(&chatCompletion)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	slog.Any("chatCompletion", chatCompletion)
	model := VertexIns.client.GenerativeModel(chatCompletion.Model)

	err = OpenAI2VerTexAI(c, chatCompletion, model)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return

	}
	return
}
