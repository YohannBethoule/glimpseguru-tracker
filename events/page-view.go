package events

import "glimpseguru-tracker/authent"

type PageViewEvent struct {
	PageURL     string `json:"page_url" binding:"required"`
	ReferrerURL string `json:"referrer_url"`
	Timestamp   string `json:"timestamp" binding:"required"`
	SessionID   string `json:"session_id"`
	authent.Identity
}

func (event PageViewEvent) Validate() bool {

}

func (event PageViewEvent) Store() error {

}
