# Go Fiber REST API Microservice

This is a simple REST API microservice built with Go, Fiber framework, GORM ORM, and Neon PostgreSQL database.

## Setup

1. Clone the repository.

2. Set up your Neon database:
   - Create a project on [Neon](https://neon.tech/).
   - Get the connection string.

3. Update the `.env` file with your database URL:
   ```
   DATABASE_URL=your_neon_connection_string_here
   ```

4. Install dependencies:
   ```
   go mod tidy
   ```

5. Run the application:
   ```
   go run main.go
   ```

The server will start on port 3000.

## Swagger Documentation

Swagger documentation is available only in development environment.

To enable Swagger in development:
1. Install the swag CLI tool: `go install github.com/swaggo/swag/cmd/swag@latest`
2. Run `swag init` to generate/update the documentation files
3. Set the environment variable `APP_ENV=development`
4. The Swagger UI will be available at `/swagger`

Note: Basic swagger.json and swagger.yaml files have been created. For production, regenerate them with `swag init` after any code changes.

## API Endpoints

- `GET /users` - Get all users
- `POST /users` - Create a new user (JSON body: {"name": "string", "email": "string"})
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user by ID
- `DELETE /users/:id` - Delete user by ID

## User Model

```go
type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```