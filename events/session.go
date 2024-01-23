package events

import (
	"context"
	"errors"
	"glimpseguru-tracker/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type EventSessionType int

const (
	Start EventSessionType = iota
	End
)

const sessionExpirationInMinutes = 10

type SessionEvent struct {
	Timestamp int64 `json:"timestamp" binding:"required"`
	Type      EventSessionType
	SessionID string `json:"session_id" binding:"required"`
	WebsiteID string `json:"website_id"`
	UserID    string `json:"user_id"`
}

func (event *SessionEvent) validate() bool {
	if event.SessionID == "" {
		return false
	}

	if event.WebsiteID == "" || event.UserID == "" {
		return false
	}

	return true
}

func (event *SessionEvent) getSessionFilter() bson.M {
	return bson.M{"sessionID": event.SessionID, "userID": event.UserID, "websiteID": event.WebsiteID}
}

func (event *SessionEvent) find() (bson.M, error) {
	var session bson.M
	errFindSession := db.MongoDb.Collection("sessions").FindOne(context.Background(), event.getSessionFilter()).Decode(&session)
	return session, errFindSession
}

func (event *SessionEvent) create() error {
	doc := bson.M{
		"start":     time.Unix(event.Timestamp, 0),
		"end":       time.Unix(event.Timestamp+(sessionExpirationInMinutes*60), 0),
		"sessionID": event.SessionID,
		"userID":    event.UserID,
		"websiteID": event.WebsiteID,
	}
	// Create session record
	_, err := db.MongoDb.Collection("sessions").InsertOne(context.Background(), doc)
	return err
}

func (event *SessionEvent) updateEndTime(endTime time.Time) error {
	filter := event.getSessionFilter()
	update := bson.M{"$set": bson.M{"end": endTime}}
	_, err := db.MongoDb.Collection("sessions").UpdateOne(context.Background(), filter, update)
	return err
}

func (event *SessionEvent) createOrUpdateSession() error {
	_, errFindSession := event.find()

	if errors.Is(errFindSession, mongo.ErrNoDocuments) {
		// Create session record
		return event.create()
	} else if errFindSession == nil {
		// Update session record
		return event.updateEndTime(time.Unix(event.Timestamp+(sessionExpirationInMinutes*60), 0))
	}

	return errFindSession
}

func (event *SessionEvent) Process(userID string, websiteID string) error {
	event.SetUser(userID, websiteID)
	if !event.validate() {
		return errors.New("invalid event")
	}

	if event.Type == Start {
		return event.createOrUpdateSession()
	} else {
		return event.updateEndTime(time.Unix(event.Timestamp, 0))
	}
}

func (event *SessionEvent) SetUser(userID string, websiteID string) {
	event.UserID = userID
	event.WebsiteID = websiteID
}
