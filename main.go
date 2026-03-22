package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
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

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Create Fiber app
	app := fiber.New()

	// Routes
	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User
		db.Find(&users)
		if db.Error != nil {
			log.Println("Error fetching users:", db.Error)
			return c.Status(500).SendString("Failed to fetch users")
		}
		return c.JSON(users)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if err := db.Create(user).Error; err != nil {
			log.Println("Error creating user:", err)
			return c.Status(500).SendString("Failed to create user")
		}
		return c.JSON(user)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			return c.Status(404).SendString("User not found")
		}
		return c.JSON(user)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			return c.Status(404).SendString("User not found")
		}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if err := db.Save(&user).Error; err != nil {
			log.Println("Error updating user:", err)
			return c.Status(500).SendString("Failed to update user")
		}
		return c.JSON(user)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			return c.Status(404).SendString("User not found")
		}
		if err := db.Delete(&user).Error; err != nil {
			log.Println("Error deleting user:", err)
			return c.Status(500).SendString("Failed to delete user")
		}
		return c.SendString("User deleted")
	})

	log.Fatal(app.Listen(":3000"))
}
