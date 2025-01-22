package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dating-app/models"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ValidateUser(username, password string) (models.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) Swipe(userID, targetID, action string) error {
	args := m.Called(userID, targetID, action)
	return args.Error(0)
}

func (m *MockUserService) PurchasePremium(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) AddVerifiedLabel(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) RemoveSwipeQuota(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) Signup(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestLoginHandler(t *testing.T) {
	mockUserService := new(MockUserService)
	userService = mockUserService

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   map[string]string
		mockReturn     struct {
			user models.User
			err  error
		}
		expectValidateUserCall bool
	}{
		{
			name: "Successful login",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Login successful"},
			mockReturn: struct {
				user models.User
				err  error
			}{
				user: models.User{Username: "testuser"},
				err:  nil,
			},
			expectValidateUserCall: true,
		},
		{
			name: "Invalid request payload",
			requestBody: map[string]string{
				"username": "testuser",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Invalid request payload"},
			mockReturn: struct {
				user models.User
				err  error
			}{
				user: models.User{},
				err:  nil,
			},
			expectValidateUserCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectValidateUserCall {
				mockUserService.On("ValidateUser", tt.requestBody["username"], tt.requestBody["password"]).Return(tt.mockReturn.user, tt.mockReturn.err)
			}

			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(LoginHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var responseBody map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
				t.Fatal(err)
			}

			for key, value := range tt.expectedBody {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tt.expectedBody)
				}
			}

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestSignupHandler(t *testing.T) {
	mockUserService := new(MockUserService)
	userService = mockUserService

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   map[string]string
		mockReturn     error
	}{
		{
			name: "Successful signup",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpassword",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   map[string]string{"message": "Signup successful"},
			mockReturn:     nil,
		},
		{
			name: "Invalid request payload",
			requestBody: map[string]string{
				"username": "testuser",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Invalid request payload"},
			mockReturn:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockReturn == nil {
				mockUserService.On("Signup", mock.Anything).Return(tt.mockReturn)
			}

			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SignupHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var responseBody map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
				t.Fatal(err)
			}

			for key, value := range tt.expectedBody {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tt.expectedBody)
				}
			}

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestSwipeHandler(t *testing.T) {
	mockUserService := new(MockUserService)
	userService = mockUserService

	tests := []struct {
		name            string
		requestBody     map[string]string
		expectedStatus  int
		expectedBody    map[string]string
		mockReturn      error
		expectSwipeCall bool
	}{
		{
			name: "Successful swipe",
			requestBody: map[string]string{
				"userID":   "user1",
				"targetID": "target1",
				"action":   "right",
			},
			expectedStatus:  http.StatusOK,
			expectedBody:    map[string]string{"message": "Swipe action recorded"},
			mockReturn:      nil,
			expectSwipeCall: true,
		},
		{
			name: "Invalid request payload",
			requestBody: map[string]string{
				"userID": "user1",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    map[string]string{"error": "Invalid request payload"},
			mockReturn:      nil,
			expectSwipeCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectSwipeCall {
				mockUserService.On("Swipe", tt.requestBody["userID"], tt.requestBody["targetID"], tt.requestBody["action"]).Return(tt.mockReturn)
			}

			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/swipe", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SwipeHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var responseBody map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
				t.Fatal(err)
			}

			for key, value := range tt.expectedBody {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tt.expectedBody)
				}
			}

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestPurchaseHandler(t *testing.T) {
	mockUserService := new(MockUserService)
	userService = mockUserService

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   map[string]string
		mockReturn     error
	}{
		{
			name: "Successful purchase - remove quota",
			requestBody: map[string]string{
				"userID":       "user1",
				"purchaseType": "remove_quota",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Purchase action completed"},
			mockReturn:     nil,
		},
		{
			name: "Successful purchase - add verified",
			requestBody: map[string]string{
				"userID":       "user1",
				"purchaseType": "add_verified",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Purchase action completed"},
			mockReturn:     nil,
		},
		{
			name: "Invalid request payload",
			requestBody: map[string]string{
				"userID": "user1",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Invalid request payload"},
			mockReturn:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockReturn == nil {
				if tt.requestBody["purchaseType"] == "remove_quota" {
					mockUserService.On("RemoveSwipeQuota", tt.requestBody["userID"]).Return(tt.mockReturn)
				} else if tt.requestBody["purchaseType"] == "add_verified" {
					mockUserService.On("AddVerifiedLabel", tt.requestBody["userID"]).Return(tt.mockReturn)
				}
			}

			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/purchase", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(PurchaseHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var responseBody map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
				t.Fatal(err)
			}

			for key, value := range tt.expectedBody {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tt.expectedBody)
				}
			}

			mockUserService.AssertExpectations(t)
		})
	}
}
