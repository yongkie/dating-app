package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Premium   bool      `json:"premium"`
	Swipes    int       `json:"swipes"`
	LastSwipe time.Time `json:"last_swipe"`
}
