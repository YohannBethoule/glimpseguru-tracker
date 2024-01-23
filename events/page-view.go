package events

import (
	"context"
	"errors"
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
	WebsiteID   string `json:"website_id"`
	UserID      string `json:"user_id"`
}

func (event *PageViewEvent) validate() bool {
	if event.PageURL == "" || event.Timestamp <= 0 {
		return false
	}

	if event.WebsiteID == "" || event.UserID == "" {
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
		"userID":      event.UserID,
		"websiteID":   event.WebsiteID,
	}

	// Insert into MongoDB collection
	_, err := db.MongoDb.Collection("page_views").InsertOne(context.Background(), doc)
	return err
}

func (event *PageViewEvent) Process(userID string, websiteID string) error {
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

func (event *PageViewEvent) SetUser(userID string, websiteID string) {
	event.UserID = userID
	event.WebsiteID = websiteID
}
