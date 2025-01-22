package service

import (
	"dating-app/models"
)

// UserService interface
type UserService interface {
	Signup(user models.User) error
	Login(username string) (*models.User, error)
	Swipe(userID, targetID, action string) error
	PurchasePremium(userID string) error
	RemoveSwipeQuota(userID string) error
	AddVerifiedLabel(userID string) error
	ValidateUser(username, password string) (models.User, error)
}
