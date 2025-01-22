package service

import (
	"database/sql"
	"dating-app/models"
)

var db *sql.DB

// UserServiceImpl struct implementing UserService
type UserServiceImpl struct{}

func (s *UserServiceImpl) Signup(user models.User) error {
	_, err := db.Exec("INSERT INTO users (id, username, premium, swipes, last_swipe) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Username, user.Premium, user.Swipes, user.LastSwipe)
	return err
}

func (s *UserServiceImpl) Login(username string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, premium, swipes, last_swipe FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Premium, &user.Swipes, &user.LastSwipe)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) Swipe(userID, targetID string) error {
	// Implement swipe logic
	return nil
}

func (s *UserServiceImpl) PurchasePremium(userID string) error {
	_, err := db.Exec("UPDATE users SET premium=true WHERE id=$1", userID)
	return err
}
