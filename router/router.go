package router

import (
	"github.com/gin-gonic/gin"
	"glimpseguru-tracker/authent"
	"glimpseguru-tracker/events"
	"net/http"
)

func processTrackingEvent(c *gin.Context, event events.Event) {
	errUser := event.SetUser(c)
	if errUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errUser.Error()})
		return
	}
	if errProcess := event.Process(); errProcess != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errProcess.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event tracked successfully"})
}

func trackPageView(c *gin.Context) {
	var event events.PageViewEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	processTrackingEvent(c, &event)
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
		apiKey := c.GetHeader("X-API-Key")
		website := c.GetHeader("X-Website-ID")

		if apiKey == "" || website == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API Key and website are required in headers"})
			return
		}

		user, errAuthent := authent.GetUser(apiKey, website)
		if errAuthent != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key or Website ID"})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
