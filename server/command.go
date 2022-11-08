package server

import "github.com/gin-gonic/gin"

type command interface {
	StreamOutput(c *gin.Context, cChan chan Log) (error, string)
}

func StreamByCommand(r command, c *gin.Context, cChan chan Log) (error, string) {
	err, data := r.StreamOutput(c, cChan)
	if err != nil {
		return err, ""
	}
	return nil, data

}
