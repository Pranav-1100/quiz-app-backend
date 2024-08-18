# Quiz App Backend

This is the backend for a gamified quiz application built with Go, Gin, and SQLite.

## Features

- Multiple quiz levels with increasing difficulty
- In-quiz currency system
- User progress tracking
- Lifeline system
- Achievements
- Leaderboard

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run the server: `go run cmd/server/main.go`

## API Endpoints

- GET /questions/:level - Get questions for a specific level
- POST /answer - Submit an answer
- GET /user/:id - Get user information
- POST /lifeline - Use a lifeline
- GET /achievements/:id - Get user achievements
- GET /leaderboard - Get the leaderboard
- GET /next-level - Get information about the next level

## Database

The application uses SQLite as its database. The schema is defined in `migrations/001_create_tables.sql`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.