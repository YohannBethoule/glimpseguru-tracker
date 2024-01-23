package events

type Event interface {
	Process(userID string, websiteID string) error
}
