package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hassanjawwad12/go-fiber-postgres/models"
	"github.com/hassanjawwad12/go-fiber-postgres/storage"
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
	books := &[]models.Books{}

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

func (r *Respository) DeleteBook(context *fiber.Ctx) error {

	// Get id from the route
	id := context.Params("id")

	// handling empty id
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}
	book := &models.Books{}

	// Delete the book with the specified id
	err := r.DB.Delete(book, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find the book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})
	return nil
}

func (r *Respository) GetBookByID(context *fiber.Ctx) error {

	// Get id from the route
	id := context.Params("id")

	// handling empty id
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}
	book := &models.Books{}

	// Find the book with the specified id
	err := r.DB.Where("id = ?", id).First(book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find the book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data":    book,
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

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	/// set up the db
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not connect to the database")
	}

	// Migrate the schema
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate the schema")
	}

	// r becomes a struct of type repository
	r := Respository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}
