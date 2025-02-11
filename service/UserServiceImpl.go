package service

import (
	"database/sql"
	"dating-app/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl struct implementing UserService
type UserServiceImpl struct {
	DB *sql.DB
}

func (s *UserServiceImpl) Signup(user models.User) error {
	_, err := s.DB.Exec("INSERT INTO users (id, username, premium, swipes, last_swipe) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Username, user.Premium, user.Swipes, user.LastSwipe)
	return err
}

func (s *UserServiceImpl) Login(username string) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRow("SELECT id, username, premium, swipes, last_swipe FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Premium, &user.Swipes, &user.LastSwipe)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) Swipe(userID, targetID, action string) error {
	var user models.User
	err := s.DB.QueryRow("SELECT swipes, last_swipe FROM users WHERE id=$1", userID).Scan(&user.Swipes, &user.LastSwipe)
	if err != nil {
		return err
	}

	// Check if the user has swiped more than 10 times today
	if user.Swipes >= 10 && user.LastSwipe.After(time.Now().AddDate(0, 0, -1)) {
		return fmt.Errorf("daily swipe limit reached")
	}

	// Check if the user has already swiped on the target profile today
	var count int
	err = s.DB.QueryRow("SELECT COUNT(*) FROM swipes WHERE user_id=$1 AND target_id=$2 AND DATE(created_at)=CURRENT_DATE", userID, targetID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("already swiped on this profile today")
	}

	// Update the user's swipe count and last swipe time
	_, err = s.DB.Exec("UPDATE users SET swipes=swipes+1, last_swipe=$1 WHERE id=$2", time.Now(), userID)
	if err != nil {
		return err
	}

	// Record the swipe action (left or right)
	_, err = s.DB.Exec("INSERT INTO swipes (user_id, target_id, action) VALUES ($1, $2, $3)", userID, targetID, action)
	return err
}

func (s *UserServiceImpl) PurchasePremium(userID string) error {
	_, err := s.DB.Exec("UPDATE users SET premium=true WHERE id=$1", userID)
	return err
}

func (s *UserServiceImpl) RemoveSwipeQuota(userID string) error {
	_, err := s.DB.Exec("UPDATE users SET swipes=0 WHERE id=$1", userID)
	return err
}

func (s *UserServiceImpl) AddVerifiedLabel(userID string) error {
	_, err := s.DB.Exec("UPDATE users SET verified=true WHERE id=$1", userID)
	return err
}

func (s *UserServiceImpl) ValidateUser(username, password string) (models.User, error) {
	var user models.User
	err := s.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Password)
	if err != nil {
		return user, fmt.Errorf("user not found")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, fmt.Errorf("invalid password")
	}

	return user, nil
}
