package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	f, _ := os.Create("app.db")
	f.Close()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		var userInput UserInput
		c.BodyParser(&userInput)
		db.Create(&User{
			Name:      userInput.Name,
			UserName:  userInput.UserName,
			Email:     userInput.Email,
			Bio:       userInput.Bio,
			BirthDate: userInput.BirthDate,
			Gender:    userInput.Gender,
			Avatar:    userInput.Avatar,
			Header:    userInput.Header,
			Password:  userInput.Password,
		})

		return c.JSON(UserInput{
			Name:      userInput.Name,
			UserName:  userInput.UserName,
			Email:     userInput.Email,
			Bio:       userInput.Bio,
			BirthDate: userInput.BirthDate,
			Gender:    userInput.Gender,
			Avatar:    userInput.Avatar,
			Header:    userInput.Header,
			Password:  userInput.Password,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var users []User

		var UsersCount int64
		db.Find(&User{}).Count(&UsersCount)
		rand := rand.Intn(int(UsersCount) - 50)
		db.Offset(rand).Limit(50).Find(&users)

		return c.JSON(users)
	})
	fmt.Println("Server Running on 5000")
	app.Listen(":5000")
}

type UserInput struct {
	Name      string `json:"name"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Header    string `json:"header"`
	Password  string `json:"password"`
}

type User struct {
	gorm.Model
	Name      string `json:"name"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Header    string `json:"header"`
	Password  string `json:"password"`
}
