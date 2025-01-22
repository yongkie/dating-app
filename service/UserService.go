package service

import (
	"dating-app/models"
)

// UserService interface
type UserService interface {
	Signup(user models.User) error
	Login(username string) (*models.User, error)
	Swipe(userID, targetID string) error
	PurchasePremium(userID string) error
}
