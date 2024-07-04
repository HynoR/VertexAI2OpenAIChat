package main

// gin
import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

var ModelConvertMap map[string]string

func init() {
	ModelConvertMap = make(map[string]string)
	ModelConvertMap["gemini-1.5-flash"] = "gemini-1.5-flash-001"
	ModelConvertMap["gemini-1.5-pro"] = "gemini-1.5-pro-001"
	ModelConvertMap["claude-3-5-sonnet-20240620"] = "publishers/anthropic/models/claude-3-5-sonnet"
	ModelConvertMap["claude-3-opus-20240229"] = "publishers/anthropic/models/claude-3-opus"
	ModelConvertMap["claude-3-haiku-20240307"] = "publishers/anthropic/models/claude-3-haiku"
	ModelConvertMap["claude-3-sonnet-20240229"] = "publishers/anthropic/models/claude-3-sonnet"
}

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
	if covertModel, ok := ModelConvertMap[chatCompletion.Model]; ok {
		chatCompletion.Model = covertModel
	}
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
