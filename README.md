# Dating App Backend

## Overview
This is a simple backend service for a dating app built with Go and PostgreSQL. It provides user signup, login, swiping, and premium purchase functionalities.

## Project Structure
```
├── main.go          # Main application entry point
├── .env             # Environment variables
├── go.mod           # Go module file
├── go.sum           # Go dependencies
├── README.md        # Documentation
```

## Prerequisites
- Go 1.16+
- PostgreSQL
- Docker (optional)

## Environment Variables
Create a `.env` file in the project root with the following:
```
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
```

## Running the Service
1. Install dependencies:
   ```sh
   go mod tidy
   ```
2. Start PostgreSQL (if not already running):
   ```sh
   docker run --name postgres -e POSTGRES_USER=your_username -e POSTGRES_PASSWORD=your_password -e POSTGRES_DB=your_database -p 5432:5432 -d postgres
   ```
3. Run the service:
   ```sh
   go run main.go
   ```
4. The server will start at `http://localhost:8080`

## API Endpoints
- **Signup** (POST `/signup`)
- **Login** (POST `/login`)
- **Swipe** (POST `/swipe`)
- **Purchase Premium** (POST `/purchase`)

## Swagger Documentation
Swagger UI is available at `/swagger/` when the server is running.

## Testing
To test the API endpoints, you can use:
```sh
curl -X POST http://localhost:8080/signup -d '{"id":"1", "username":"testuser", "premium":false, "swipes":10, "last_swipe":"2023-01-01T00:00:00Z"}' -H "Content-Type: application/json"
```
