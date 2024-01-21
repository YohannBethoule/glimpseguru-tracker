package events

import "github.com/gin-gonic/gin"

type Event interface {
	SetUser(c *gin.Context) error
	Process() error
}
