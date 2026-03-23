// Package main GoFiber REST API Microservice
//
// @title       GoFiber REST API
// @version     1.0
// @description A simple REST API built with Go Fiber
// @host        localhost:3000
// @BasePath    /
package main

import (
	"log"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" gorm:"unique" example:"john@example.com"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db

	// Migrate the schema
	DB.AutoMigrate(&User{})

	// Create Fiber app
	app := fiber.New()

	// Enable Swagger only in development
	// Note: Run 'swag init -g main.go -o docs' to generate docs/swagger.json and docs/swagger.yaml
	if os.Getenv("APP_ENV") == "development" {
		app.Use(swagger.New(swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "swagger",
			Title:    "GoFiber REST API Docs",
		}))
	}

	app.Get("/users", GetUsers)
	app.Post("/users", CreateUser)
	app.Get("/users/:id", GetUserByID)
	app.Put("/users/:id", UpdateUser)
	app.Delete("/users/:id", DeleteUser)

	log.Fatal(app.Listen(":3000"))
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users from the database
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Failure 500 {string} string "Failed to fetch users"
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		return c.Status(500).SendString("Failed to fetch users")
	}
	return c.JSON(users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 200 {object} User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Failed to create user"
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := DB.Create(user).Error; err != nil {
		log.Println("Error creating user:", err)
		return c.Status(500).SendString("Failed to create user")
	}
	return c.JSON(user)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body User true "User data"
// @Success 200 {object} User
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Failed to update user"
// @Router /users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := DB.Save(&user).Error; err != nil {
		log.Println("Error updating user:", err)
		return c.Status(500).SendString("Failed to update user")
	}
	return c.JSON(user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Failed to delete user"
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	if err := DB.Delete(&user).Error; err != nil {
		log.Println("Error deleting user:", err)
		return c.Status(500).SendString("Failed to delete user")
	}
	return c.SendString("User deleted")
}
