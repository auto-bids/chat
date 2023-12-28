package Controllers

import "github.com/gin-gonic/gin"

func GetMessages(ctx *gin.Context) {
	go func(c *gin.Context) {
		//todo - baza
	}(ctx.Copy())
}
