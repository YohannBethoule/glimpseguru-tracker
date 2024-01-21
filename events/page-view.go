package events

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"glimpseguru-tracker/authent"
	"glimpseguru-tracker/db"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"time"
)

type PageViewEvent struct {
	PageURL     string `json:"page_url" binding:"required"`
	Pathname    string `json:"pathname" binding:"required"`
	ReferrerURL string `json:"referrer_url"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
	SessionID   string `json:"session_id"`
	DeviceType  string `json:"device_type"`
	SourceType  string `json:"source_type"`
	*authent.User
}

func (event *PageViewEvent) validate() bool {
	if event.PageURL == "" || event.Timestamp <= 0 {
		return false
	}

	if event.User == nil {
		return false
	}

	return true
}

func (event *PageViewEvent) store() error {
	doc := bson.M{
		"pageURL":     event.PageURL,
		"pathname":    event.Pathname,
		"referrerURL": event.ReferrerURL,
		"timestamp":   time.Unix(event.Timestamp, 0),
		"sessionID":   event.SessionID,
		"deviceType":  event.DeviceType,
		"sourceType":  event.SourceType,
		"user":        event.User,
	}

	// Insert into MongoDB collection
	_, err := db.MongoDb.Collection("page_views").InsertOne(context.Background(), doc)
	return err
}

func (event *PageViewEvent) Process() error {
	if !event.validate() {
		return errors.New("invalid event")
	}

	return event.store()
}

func (event *PageViewEvent) SetUser(c *gin.Context) error {
	if storedUser, exists := c.Get("user"); exists {
		// Type assert to ensure it's of the Identity type
		if user, ok := storedUser.(authent.User); ok {
			event.User = &user
		} else {
			return errors.New("invalid user identity")
		}
	} else {
		slog.Warn("no user identity found in context")
		return errors.New("user identity required")
	}
	return nil
}
