package router

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"glimpseguru-tracker/authent"
	"glimpseguru-tracker/events"
	"net/http"
	"time"
)

func processTrackingEvent(c *gin.Context, event events.Event) {
	var errUser error
	var user authent.User
	var isUserStruct bool
	if storedUser, exists := c.Get("user"); exists {
		// Type assert to ensure it's of the Identity type
		if user, isUserStruct = storedUser.(authent.User); !isUserStruct {
			errUser = errors.New("cannot bind user from context")
		}
	} else {
		errUser = errors.New("no user found in context")
	}
	if errUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error setting user: %e", errUser)})
		return
	}

	if errProcess := event.Process(user.ID, user.Website.ID); errProcess != nil {
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

func trackSessionEnd(c *gin.Context) {
	var event events.SessionEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.Type = events.End
	processTrackingEvent(c, &event)
}

func trackCustomEvent(c *gin.Context) {
	var event events.CustomEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	processTrackingEvent(c, &event)
}

func New() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST", "OPTIONS"},
		AllowHeaders:    []string{"x-api-key", "x-website-id", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))
	rTracking := r.Group("/track/")
	rTracking.Use(identityValidationMiddleware())
	rTracking.POST("/pageview", trackPageView)
	rTracking.POST("/sessionend", trackSessionEnd)
	rTracking.POST("/customevent", trackCustomEvent)

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
