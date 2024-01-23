package events

import (
	"glimpseguru-tracker/authent"
)

type Event interface {
	SetUser(user authent.User) error
	Process() error
}
