package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/vercel/go-middleware"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books = []Book{
    {ID: 1, Title: "The Catcher in the Rye", Author: "J.D. Salinger"},
    {ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"},
    {ID: 3, Title: "1984", Author: "George Orwell"},
}

func main() {
    app := fiber.New()

    // Wrap the app in a Vercel-compatible middleware handler
    handler := middleware.New(app)

    // GET all books
    app.Get("/api/books", func(c *fiber.Ctx) error {
        return c.JSON(books)
    })

    // GET a book by ID
    app.Get("/api/books/:id", func(c *fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.Status(http.StatusBadRequest).SendString("Bad request")
        }
        for _, book := range books {
            if book.ID == id {
                return c.JSON(book)
            }
        }
        return c.Status(http.StatusNotFound).SendString("Book not found")
    })

    // POST a new book
    app.Post("/api/books", func(c *fiber.Ctx) error {
        var book Book
        if err := json.Unmarshal(c.Body(), &book); err != nil {
            return c.Status(http.StatusBadRequest).SendString("Bad request")
        }
        books = append(books, book)
        return c.SendString("Book added")
    })

    // DELETE a book by ID
    app.Delete("/api/books/:id", func(c *fiber.Ctx) error {
        id, err := strconv.Atoi(c.Params("id"))
        if err != nil {
            return c.Status(http.StatusBadRequest).SendString("Bad request")
        }
        for i, book := range books {
            if book.ID == id {
                books = append(books[:i], books[i+1:]...)
                return c.SendString("Book deleted")
            }
        }
        return c.Status(http.StatusNotFound).SendString("Book not found")
    })

    // Start server
    log.Fatal(http.ListenAndServe("", handler))
}