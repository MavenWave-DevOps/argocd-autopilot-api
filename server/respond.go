package server

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int
	Message string `json:"message"`
	Error   string `json:"error"`
	Link    string `json:"link"`
}

func SendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, gin.H{"message": response.Message, "link": response.Link})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, response.Error)
	}
}
