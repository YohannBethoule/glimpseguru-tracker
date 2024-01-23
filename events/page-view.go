package events

import (
	"context"
	"errors"
	"glimpseguru-tracker/authent"
	"glimpseguru-tracker/db"
	"go.mongodb.org/mongo-driver/bson"
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
		"user_id":     event.UserID,
		"website_id":  event.WebsiteID,
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

func (event *PageViewEvent) SetUser(user authent.User) error {
	event.UserID = user.ID
	event.WebsiteID = user.Website.ID
	return nil
}
