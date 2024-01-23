package authent

import (
	"errors"
	"fmt"
	"glimpseguru-tracker/db"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"log/slog"
)

type Identity struct {
	ApiKey  string `json:"api_key" binding:"required"`
	Website string `json:"website" binding:"required"`
}

type User struct {
	ID string `bson:"userID"`
	Website
}

type UserResponse struct {
	ID       string    `bson:"userID"`
	APIKey   string    `bson:"apiKey"`
	Websites []Website `bson:"websites"`
}

type Website struct {
	ID   string `bson:"websiteID"`
	Name string `bson:"websiteName"`
	URL  string `bson:"websiteURL"`
}

func GetUser(apiKey string, websiteID string) (User, error) {
	var userResult UserResponse
	var user User

	// Get user
	filter := bson.M{
		"apiKey":             apiKey,
		"websites.websiteID": websiteID,
	}
	result := db.MongoDb.Collection("users").FindOne(context.Background(), filter)
	if result.Err() != nil {
		slog.Error("no user found for apiKey and website", slog.String("api key", apiKey), slog.String("website", websiteID), slog.String("error", result.Err().Error()))
		return user, errors.New("no user found for apiKey and website")
	}
	err := result.Decode(&userResult)
	if err != nil {
		slog.Error(fmt.Sprintf("error decoding user: %e", err))
		return user, errors.New("error decoding user response")
	}

	// Get asked website
	for _, website := range userResult.Websites {
		if website.ID == websiteID {
			user = User{
				ID: userResult.ID,
				Website: Website{
					ID:   website.ID,
					Name: website.Name,
					URL:  website.URL,
				},
			}
		}
	}

	return user, nil
}
