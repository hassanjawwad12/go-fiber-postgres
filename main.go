package main

import (
	"log"
	"net/http"

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

func (r *Respository) CreateBook(context *fiber.Ctx) error {

	book := Book{}

	// Convert JSON to Book format
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "could not parse the request",
		})
		return err
	}

	// Add it to the database
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create the book",
		})
		return err
	}
	context.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "book created successfully",
	})

	// no error so we return nil
	return nil
}

func (r *Respository) GetBooks(context *fiber.Ctx) error {
	books := &[]models.Book{}

	err := r.DB.Find(books).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not fetch the books",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    books,
	})
	return nil
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
