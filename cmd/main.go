package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/snehabhatia04/libmgmt/controllers"
	"github.com/snehabhatia04/libmgmt/web"
)

func main() {
	// Database connection
	db, err := sqlx.Connect("postgres", "postgres://postgres:Sneha@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	// Initializing stores
	userStore := controllers.NewDBUserStore(db)
	bookStore := controllers.NewDBBookStore(db)
	IssuedBookStore := controllers.NewDBIssuedBookStore(db)
	FineStore := controllers.NewDBFineStore(db)

	// Initializing handlers
	handler := web.NewHandler(userStore, bookStore, IssuedBookStore, FineStore)

	// Setting up router
	router := gin.Default()

	// User routes
	router.GET("/users", handler.GetUsers)
	router.GET("/users/:id", handler.GetUser)
	router.POST("/users", handler.CreateUser)
	router.DELETE("/users/:id", handler.DeleteUser)

	// Book routes
	router.GET("/books", handler.GetBooks)
	router.GET("/books/:id", handler.GetBook)
	router.POST("/books", handler.CreateBook)
	router.PUT("/books/:id", handler.UpdateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	// Issued Book routes
	router.POST("/books/issue", handler.IssueBook)
	router.PUT("/books/return/:id", handler.ReturnBook)
	// router.GET("/books/issued", handler.GetIssuedBooks)
	// router.GET("/books/issued/:id", handler.GetIssuedBook)

	// Fine routes
	router.GET("/fines", handler.GetFines)
	router.GET("/fines/:id", handler.GetFine)
	router.DELETE("/fines/:id", handler.DeleteFine)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start the server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
