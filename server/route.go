package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Hostname string
	Port     int
}

type Stringer interface {
	Start() error
}

func (r Server) Start() error {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/run", ExecuteCommands)
	if err := router.Run(fmt.Sprintf("%s:%s", r.Hostname, r.ToString())); err != nil {
		return err
	}
	return nil
}
