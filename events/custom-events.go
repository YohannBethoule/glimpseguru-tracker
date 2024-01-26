package events

import (
	"context"
	"errors"
	"glimpseguru-tracker/db"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"time"
)

type CustomEvent struct {
	Label     string `json:"label" binding:"required"`
	Data      any    `json:"data"`
	Timestamp int64  `json:"timestamp" binding:"required"`
	SessionID string `json:"session_id"`
	WebsiteID string `json:"website_id"`
	UserID    string `json:"user_id"`
}

func (event *CustomEvent) validate() bool {
	if event.Label == "" || event.Timestamp <= 0 {
		return false
	}

	if event.WebsiteID == "" || event.UserID == "" {
		return false
	}

	return true
}

func (event *CustomEvent) store() error {
	doc := bson.M{
		"label":     event.Label,
		"data":      event.Data,
		"timestamp": time.Unix(event.Timestamp, 0),
		"sessionID": event.SessionID,
		"userID":    event.UserID,
		"websiteID": event.WebsiteID,
	}

	// Insert into MongoDB collection
	_, err := db.MongoDb.Collection("custom_events").InsertOne(context.Background(), doc)
	return err
}

func (event *CustomEvent) Process(userID string, websiteID string) error {
	event.SetUser(userID, websiteID)
	if !event.validate() {
		return errors.New("invalid event")
	}

	errStore := event.store()
	if errStore != nil {
		return errStore
	}

	session := SessionEvent{
		Timestamp: event.Timestamp,
		Type:      Start,
		SessionID: event.SessionID,
	}
	errSession := session.Process(userID, websiteID)
	if errSession != nil {
		slog.Error("error processing session", slog.String("error", errSession.Error()))
	}
	return nil
}

func (event *CustomEvent) SetUser(userID string, websiteID string) {
	event.UserID = userID
	event.WebsiteID = websiteID
}
