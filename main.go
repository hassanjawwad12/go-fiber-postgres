package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// Created our own data type which have the db
type Respository struct {
	DB *gorm.DB
}

// SetupRoutes is a struct method
func (r *Respository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// All the api belong to the api group
	// We are calling method here not functions
	api.Post("/create_book", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	/// set up the db
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not connect to the database")
	}

	// r becomes a struct of type repository
	r := Respository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}
