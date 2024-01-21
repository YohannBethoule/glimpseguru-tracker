package router

import (
	"github.com/gin-gonic/gin"
	"glimpseguru-tracker/events"
	"glimpseguru-tracker/identity"
	"net/http"
)

func trackPageView(c *gin.Context) {
	var event events.PageViewEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process and store the event in MongoDB
	c.JSON(http.StatusOK, gin.H{"message": "Page view tracked successfully"})
}

func New() *gin.Engine {
	r := gin.Default()

	rTracking := r.Group("/track/")
	rTracking.Use(identityValidationMiddleware())
	rTracking.POST("/pageview", trackPageView)

	return r
}

func identityValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var identityQuery identity.Identity
		if err := c.ShouldBindJSON(&identityQuery); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if !identity.Validate(identityQuery) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key or Website ID"})
			return
		}

		// Set identity for downstream use
		c.Set("identity", identityQuery)

		c.Next()
	}
}
